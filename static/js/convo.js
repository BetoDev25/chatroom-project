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
    if (!tabBar) {
        console.error('convoTabBar not found');
        return;
    }

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
    
    loadConvoMessages(convoId, true);
    saveConvoTabs();
}

function switchConvoTab(convoId) {
    const tab = convoTabs.find(t => t.ConversationID === convoId);
    if (!tab) return;

    activeConvoTabId = convoId;
    currentConversation = tab;

    document.querySelectorAll('.convo-tab').forEach(t => t.classList.remove('active'));
    document.querySelector(`.convo-tab[data-convo-id="${convoId}"]`).classList.add('active');

    loadConvoMessages(convoId, true);
    saveConvoTabs();
}

function closeConvoTab(convoId) {
    const tabElement = document.querySelector(`.convo-tab[data-convo-id="${convoId}"]`);
    if (tabElement) tabElement.remove();

    convoTabs = convoTabs.filter(t => t.ConversationID !== convoId);

    if (activeConvoTabId === convoId) {
        if (convoTabs.length > 0) {
            switchConvoTab(convoTabs[0].ConversationID);
        } else {
            activeConvoTabId = null;
            currentConversation = null;
            document.getElementById('convoMessages').innerHTML = '';
        }
    }
    saveConvoTabs();
}

async function loadConvoMessages(convoId, resetPage = true) {
    if (resetPage) {
        convoPage = 1;
        document.getElementById('convoMessages').innerHTML = '';
    }

    const response = await fetch(`/api/priv-messages/${convoId}?page=${convoPage}&limit=${convoMessagesPerPage}`, {
        credentials: 'include'
    });
    
    if (!response.ok) {
        console.error('Failed to load conversation messages');
        return;
    }

    const messages = await response.json();
    const messagesDiv = document.getElementById('convoMessages');
    
    messages.forEach(msg => {
        const div = document.createElement('div');
        const username = msg.Username || msg.username;
        const content = msg.Content || msg.content;
        div.textContent = `${username}: ${content}`;
        messagesDiv.appendChild(div);
    });
    
    messagesDiv.scrollTop = messagesDiv.scrollHeight;
    convoPage++;
}

function sendConvoMessage() {
    const input = document.getElementById('convoMessageInput');
    const content = input.value;
    if (content == "" || !currentConversation) {
        return;
    }
    
    ws.send(JSON.stringify({
        type: 'private',
        content: content,
        conversation_id: currentConversation.ConversationID
    }));
    
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
        console.error('Failed to archive private message:', err);
    });
}

function handlePrivateMessage(msg) {
    const messagesDiv = document.getElementById('convoMessages');
    if (!messagesDiv) return;
    
    const div = document.createElement('div');
    const username = msg.username || msg.Username;
    const content = msg.content || msg.Content;
    div.textContent = `${username}: ${content}`;
    messagesDiv.appendChild(div);
    messagesDiv.scrollTop = messagesDiv.scrollHeight;
}

function saveConvoTabs() {
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
        switchConvoTab(savedActive);
    }
}

document.addEventListener('DOMContentLoaded', function() {
    loadConvoTabs();
    
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

window.sendConvoMessage = sendConvoMessage;
window.handlePrivateMessage = handlePrivateMessage;