// convo.js - completely independent from index.html

let convoTabs = [];
let activeConvoTabId = null;
let currentConversation = null;
let convoPage = 1;
const convoMessagesPerPage = 50;

window.openConversation = async function(friend) {
    try {
        const response = await fetch(`/api/conversations/${friend.FriendshipID}`, {
            credentials: 'include'
        });

        let conversation;
        if (response.status === 404) {
            const createRes = await fetch(`/api/conversations/${friend.FriendshipID}`, {
                method: 'POST',
                credentials: 'include'
            });
            conversation = await createRes.json();
        } else {
            conversation = await response.json();
        }

        if (ws && ws.readyState === WebSocket.OPEN) {
            ws.send(JSON.stringify({
                type: 'private_join',
                conversation_id: conversation.ConversationID
            }));
        } else {
            console.log('openConversation: WebSocket not open, readyState:', ws ? ws.readyState : 'ws is null'); // Debug
        }
        
        addConvoTab(conversation, friend.Username);
    } catch (error) {
        console.error('Error opening conversation:', error);
        alert('Could not open conversation');
    }
};

function addConvoTab(conversation, displayName) {
    const convoId = conversation.ConversationID;
    
    if (convoTabs.find(t => t.ConversationID === convoId)) {
        switchConvoTab(convoId);
        return;
    }
    
    const tabData = {
        ConversationID: convoId,
        FriendshipID: conversation.FriendshipID,
        DisplayName: displayName
    };
    convoTabs.push(tabData);
    activeConvoTabId = convoId;
    currentConversation = tabData;
    
    const tabBar = document.getElementById('convoTabBar');
    const tab = document.createElement('button');
    tab.className = 'convo-tab active';
    tab.dataset.convoId = convoId;
    tab.innerHTML = `${displayName} <span class="close" data-convo-id="${convoId}">✕</span>`;
    
    tab.addEventListener('click', (e) => {
        if (e.target.classList.contains('close')) {
            e.stopPropagation();
            closeConvoTab(convoId);
            return;
        }
        switchConvoTab(convoId);
    });
    
    tabBar.appendChild(tab);
    
    // Clear and load messages into #messages
    document.getElementById('messages').innerHTML = '';
    loadConvoMessages(convoId, true);
    saveConvoTabs();
}

function switchConvoTab(convoId) {
    console.log('Switching to convo tab:', convoId);
    const tab = convoTabs.find(t => t.ConversationID === convoId);
    if (!tab) return;

    activeConvoTabId = convoId;
    currentConversation = tab;
    window.currentRoom = null;

    // CLEAR room active state
    activeTabId = null;
    localStorage.setItem('activeTabId', null);

    document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
    document.querySelectorAll('.convo-tab').forEach(t => t.classList.remove('active'));
    document.querySelector(`.convo-tab[data-convo-id="${convoId}"]`).classList.add('active');

    if (ws && ws.readyState === WebSocket.OPEN) {
        console.log('switchConvoTab: Joining conversation:', convoId);
        ws.send(JSON.stringify({
            type: 'private_join',
            conversation_id: convoId
        }));
    }

    loadConvoMessages(convoId, true);
    saveConvoTabs();
}


function closeConvoTab(convoId) {
    // Leave the conversation via WebSocket
    if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({
            type: 'private_leave',
            conversation_id: convoId
        }));
    }
    
    const tabElement = document.querySelector(`.convo-tab[data-convo-id="${convoId}"]`);
    if (tabElement) tabElement.remove();

    convoTabs = convoTabs.filter(t => t.ConversationID !== convoId);

    if (activeConvoTabId === convoId) {
        if (convoTabs.length > 0) {
            switchConvoTab(convoTabs[0].ConversationID);
        } else {
            // No more convo tabs - clear everything
            activeConvoTabId = null;
            currentConversation = null;
            window.currentConversation = null;
            document.getElementById('messages').innerHTML = '';
            
            // Clear active convo state from localStorage
            localStorage.setItem('activeConvoTabId', null);
            
            // Check if there are room tabs to switch to
            if (tabs.length > 0) {
                switchTab(tabs[0].RoomID, false);
            }
        }
    }
    saveConvoTabs();
}

async function loadConvoMessages(convoId, resetPage = true) {
    if (resetPage) {
        convoPage = 1;
        document.getElementById('messages').innerHTML = '';
    }

    const response = await fetch(`/api/priv-messages/${convoId}?page=${convoPage}&limit=${convoMessagesPerPage}`, {
        credentials: 'include'
    });
    
    if (!response.ok) {
        console.error('Failed to load conversation messages');
        return;
    }

    const messages = await response.json();
    const messagesDiv = document.getElementById('messages');
    
    messages.forEach(msg => {
        const div = document.createElement('div');
        const username = msg.Username || msg.username;
        const content = msg.EncryptedContent;
        
        if (username === currentUser.username) {
            div.textContent = `${username}: ${content}`;
        } else {
            const usernameSpan = document.createElement('span');
            usernameSpan.textContent = username;
            usernameSpan.style.cssText = 'cursor: pointer; color: #0066cc; text-decoration: underline;';
            usernameSpan.className = 'chat-username';
            usernameSpan.dataset.username = username;
            usernameSpan.onclick = function(e) {
                e.stopPropagation();
                window.toggleUserDropdown(e, this.dataset.username);
            };

            const contentSpan = document.createElement('span');
            contentSpan.textContent = `: ${content}`;

            div.appendChild(usernameSpan);
            div.appendChild(contentSpan);
        }
        
        messagesDiv.appendChild(div);
    });
    
    messagesDiv.scrollTop = messagesDiv.scrollHeight;
    convoPage++;
}

function sendConvoMessage() {
    const input = document.getElementById('message');
    const content = input.value;
    
    if (content == "" || !currentConversation) {
        return;
    }
    
    console.log('ws readyState:', ws ? ws.readyState : 'ws is null'); // Debug
    
    if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({
            type: 'private',
            content: content,
            conversation_id: currentConversation.ConversationID
        }));
    } else {
        console.error('WebSocket is not open!'); // Debug
    }
    
    input.value = '';
    
    fetch('/api/priv-messages', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({
            conversation_id: currentConversation.ConversationID,
            user_id: currentUser.id,
            content: content
        })
    }).catch(err => {
        console.error('Failed to archive:', err);
    });
}

function handlePrivateMessage(msg) {
    const messagesDiv = document.getElementById('messages');
    if (!messagesDiv) return;
    
    const div = document.createElement('div');
    div.textContent = `${msg.username}: ${msg.content}`;
    messagesDiv.appendChild(div);
    messagesDiv.scrollTop = messagesDiv.scrollHeight;
}

function saveConvoTabs() {
    console.log('Saving convo tabs:', convoTabs); // Debug
    console.log('Saving activeConvoTabId:', activeConvoTabId); // Debug
    localStorage.setItem('convoTabs', JSON.stringify(convoTabs));
    localStorage.setItem('activeConvoTabId', activeConvoTabId);
}

function loadConvoTabs() {
    const savedTabs = localStorage.getItem('convoTabs');
    if (!savedTabs) return;
    
    const parsed = JSON.parse(savedTabs);
    
    parsed.forEach(tabData => {
        if (convoTabs.find(t => t.ConversationID === tabData.ConversationID)) {
            return;
        }
        
        convoTabs.push(tabData);
        
        const tabBar = document.getElementById('convoTabBar');
        if (!tabBar) return;
        
        const tab = document.createElement('button');
        tab.className = 'convo-tab';
        tab.dataset.convoId = tabData.ConversationID;
        tab.innerHTML = `${tabData.DisplayName} <span class="close" data-convo-id="${tabData.ConversationID}">✕</span>`;
        
        tab.addEventListener('click', (e) => {
            if (e.target.classList.contains('close')) {
                e.stopPropagation();
                closeConvoTab(tabData.ConversationID);
                return;
            }
            switchConvoTab(tabData.ConversationID);
        });
        
        tabBar.appendChild(tab);
    });
    
    const savedActive = localStorage.getItem('activeConvoTabId');
    if (savedActive && convoTabs.find(t => t.ConversationID === savedActive)) {
        window.currentRoom = null;
        document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));

        // Join the conversation when restoring
        if (ws && ws.readyState === WebSocket.OPEN) {
            ws.send(JSON.stringify({
                type: 'private_join',
                conversation_id: savedActive
            }));
        }
        switchConvoTab(savedActive);
    }
}

document.addEventListener('DOMContentLoaded', function() {
    // Scroll to load more messages
    const messagesDiv = document.getElementById('convoMessages');
    if (messagesDiv) {
        messagesDiv.addEventListener('scroll', function() {
            if (this.scrollTop === 0 && currentConversation) {
                loadConvoMessages(currentConversation.ConversationID, false);
            }
        });
    }
    
    // Enter key for sending
    const input = document.getElementById('convoMessageInput');
    if (input) {
        input.addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                sendConvoMessage();
            }
        });
    }
});

function addConvoTabSilently(conversationId, displayName) {
    if (convoTabs.find(t => t.ConversationID === conversationId)) {
        return;
    }
    
    const tabData = {
        ConversationID: conversationId,
        DisplayName: displayName
    };
    convoTabs.push(tabData);
    
    const tabBar = document.getElementById('convoTabBar');
    const tab = document.createElement('button');
    tab.className = 'convo-tab';
    tab.dataset.convoId = conversationId;
    tab.innerHTML = `${displayName} <span class="close" data-convo-id="${conversationId}">✕</span>`;
    
    tab.addEventListener('click', (e) => {
        if (e.target.classList.contains('close')) {
            e.stopPropagation();
            closeConvoTab(conversationId);
            return;
        }
        switchConvoTab(conversationId);
    });
    
    tabBar.appendChild(tab);
    saveConvoTabs();
}

if (typeof loadFriends === 'function') {
    loadFriends();
}
window.sendConvoMessage = sendConvoMessage;
window.handlePrivateMessage = handlePrivateMessage;