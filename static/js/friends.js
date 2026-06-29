async function loadRooms() {
    try {
        const response = await fetch('/api/rooms', {
            credentials: 'include'
        });

        if (!response.ok) {
            console.error('Failed to load rooms');
            return;
        }

        const rooms = await response.json();
        
        const roomList = document.getElementById('roomList');
        if (!roomList) {
            console.error('roomList element not found');
            return;
        }
        
        roomList.innerHTML = '';

        if (rooms.length === 0) {
            roomList.innerHTML = '<div style="color: #888; padding: 20px 0; text-align: center;">No rooms created yet.</div>';
            return;
        }

        rooms.forEach(room => {
            const div = document.createElement('div');
            div.style.cssText = 'padding: 8px 0; border-bottom: 1px solid #ddd; display: flex; align-items: center; justify-content: space-between;';

            const nameSpan = document.createElement('span');
            nameSpan.textContent = room.RoomName;

            const buttonContainer = document.createElement('div');
            buttonContainer.style.cssText = 'display: flex; align-items: center; gap: 6px;';

            const joinButton = document.createElement('button');
            joinButton.textContent = 'Join';
            joinButton.style.cssText = 'padding: 2px 12px; cursor: pointer;';
            joinButton.onclick = () => {
                const roomInput = document.getElementById('room');
                if (roomInput) {
                    roomInput.value = room.RoomName;
                    joinRoom(false);
                }
            };

            const deleteButton = document.createElement('button');
            deleteButton.style.cssText = 'padding: 2px 6px; cursor: pointer; background: none; border: 1px solid #ccc; border-radius: 3px;';
            deleteButton.onclick = async () => {
                if (confirm(`Delete room "${room.RoomName}"?`)) {
                    try {
                        const deleteResponse = await fetch(`/api/rooms/${encodeURIComponent(room.RoomName)}`, {
                            method: 'DELETE',
                            credentials: 'include'
                        });
                        
                        if (deleteResponse.ok) {
                            loadRooms(); // Refresh the list
                        } else {
                            console.error('Failed to delete room');
                        }
                    } catch (error) {
                        console.error('Error deleting room:', error);
                    }
                }
            };

            const trashIcon = document.createElement('img');
            trashIcon.src = '/static/assets/trash.png';
            trashIcon.alt = 'Delete';
            trashIcon.style.cssText = 'width: 16px; height: 16px; display: block;';
            deleteButton.appendChild(trashIcon);

            buttonContainer.appendChild(joinButton);
            buttonContainer.appendChild(deleteButton);

            div.appendChild(nameSpan);
            div.appendChild(buttonContainer);
            roomList.appendChild(div);
        });
    } catch (error) {
        console.error('Error loading rooms:', error);
    }
}

// friends.html

let currentTargetUser = null;

window.toggleUserDropdown = function(event, username) {
    const dropdown = document.getElementById('userDropdown');
    const existing = document.querySelector('.user-dropdown-open');
    
    if (existing && existing.dataset.username === username) {
        dropdown.style.display = 'none';
        document.querySelectorAll('.user-dropdown-open').forEach(el => {
            el.classList.remove('user-dropdown-open');
        });
        currentTargetUser = null;
        return;
    }
    
    document.querySelectorAll('.user-dropdown-open').forEach(el => {
        el.classList.remove('user-dropdown-open');
    });
    
    const usernameElement = event.target;
    const rect = usernameElement.getBoundingClientRect();
    
    let top = rect.bottom + 5;
    let left = rect.left;
    
    const dropdownWidth = 120;
    if (left + dropdownWidth > window.innerWidth) {
        left = window.innerWidth - dropdownWidth - 10;
    }
    
    if (left < 10) {
        left = 10;
    }
    
    if (top + 80 > window.innerHeight) {
        top = rect.top - 80;
    }
    
    dropdown.style.display = 'block';
    dropdown.style.left = left + 'px';
    dropdown.style.top = top + 'px';
    
    currentTargetUser = username;
    
    document.querySelectorAll('.chat-username').forEach(el => {
        if (el.dataset.username === username) {
            el.classList.add('user-dropdown-open');
        }
    });
};

window.handleUserAction = function(action) {
    const dropdown = document.getElementById('userDropdown');
    dropdown.style.display = 'none';
    document.querySelectorAll('.user-dropdown-open').forEach(el => {
        el.classList.remove('user-dropdown-open');
    });
    
    if (action === 'add') {
        const targetUser = currentTargetUser;
        
        if (!currentUser || !currentUser.id) {
            alert('You must be logged in to add friends');
            currentTargetUser = null;
            return;
        }
        
        fetch(`/api/users/${encodeURIComponent(targetUser)}`, {
            credentials: 'include'
        })
        .then(res => {
            if (!res.ok) throw new Error('User not found');
            return res.json();
        })
        .then(user => {
            return fetch('/api/friend-request', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
                body: JSON.stringify({
                    sender_id: currentUser.id,
                    receiver_id: user.ID
                })
            });
        })
        .then(res => {
            if (res.ok) {
                alert(`Friend request sent to ${targetUser}`);
            } else {
                alert('Failed to send friend request');
            }
        })
        .catch(err => {
            console.error('Error:', err);
            alert('Error sending friend request');
        });
        
    } else if (action === 'block') {
        alert(`Blocked ${currentTargetUser}`);
    }
    
    currentTargetUser = null;
};

// Close dropdown when clicking outside
document.addEventListener('click', function(event) {
    const dropdown = document.getElementById('userDropdown');
    const isClickInside = dropdown.contains(event.target) || 
                          event.target.classList.contains('chat-username') ||
                          event.target.closest('.chat-username');
    
    if (!isClickInside && dropdown.style.display === 'block') {
        dropdown.style.display = 'none';
        document.querySelectorAll('.user-dropdown-open').forEach(el => {
            el.classList.remove('user-dropdown-open');
        });
        currentTargetUser = null;
    }
});

async function loadFriends() {
    try {
        const response = await fetch('/api/friend-request/accepted', {
            credentials: 'include'
        });
        
        if (!response.ok) {
            console.error('Failed to load friends');
            return;
        }
        
        const friends = await response.json() || [];
        const friendList = document.getElementById('friendList');
        friendList.innerHTML = '';
        
        if (friends.length === 0) {
            friendList.innerHTML = '<div style="color: #888; padding: 20px 0; text-align: center;">nobody here. try making new friends!</div>';
            return;
        }
        
        friends.forEach(friend => {
            const div = document.createElement('div');
            div.style.cssText = 'padding: 8px 0; border-bottom: 1px solid #ddd; display: flex; align-items: center; gap: 10px;';
            
            const nameSpan = document.createElement('span');
            nameSpan.textContent = friend.Username;
            nameSpan.style.cssText = 'cursor: pointer; color: #0066cc; text-decoration: underline; flex: 1;';
            nameSpan.className = 'friend-username';
            nameSpan.dataset.friendshipId = friend.FriendshipID;
            nameSpan.dataset.username = friend.Username;
            nameSpan.onclick = function(e) {
                e.stopPropagation();
                toggleFriendDropdown(e, this.dataset.friendshipId, this.dataset.username);
            };
            
            const chatButton = document.createElement('button');
            chatButton.style.cssText = 'padding: 2px 8px; cursor: pointer; background: none; border: none;';
            chatButton.onclick = async () => {
                if (typeof openConversation === 'function') {
                    openConversation(friend);
                } else {
                    alert('Conversation system not loaded yet');
                }
            };
            
            const chatImg = document.createElement('img');
            chatImg.src = '/static/assets/chat.png';
            chatImg.alt = 'chat';
            chatImg.style.cssText = 'width: 18px; height: 18px; display: block;';
            chatButton.appendChild(chatImg);
            
            div.appendChild(nameSpan);
            div.appendChild(chatButton);
            friendList.appendChild(div);
        });
    } catch (error) {
        console.error('Error loading friends:', error);
    }
}

function toggleFriendDropdown(event, friendshipId, username) {
    const existing = document.querySelector('.friend-dropdown-open');
    if (existing) {
        existing.remove();
    }
    
    const dropdown = document.createElement('div');
    dropdown.className = 'friend-dropdown-open';
    dropdown.style.cssText = 'position: fixed; background: white; border: 1px solid #ccc; border-radius: 4px; padding: 5px 0; min-width: 120px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); z-index: 2000;';
    
    const rect = event.target.getBoundingClientRect();
    let top = rect.bottom + 5;
    let left = rect.left;
    
    if (left + 120 > window.innerWidth) {
        left = window.innerWidth - 130;
    }
    if (left < 10) {
        left = 10;
    }
    if (top + 40 > window.innerHeight) {
        top = rect.top - 40;
    }
    
    dropdown.style.left = left + 'px';
    dropdown.style.top = top + 'px';
    
    const removeOption = document.createElement('div');
    removeOption.textContent = 'Remove Friend';
    removeOption.style.cssText = 'padding: 8px 15px; cursor: pointer;';
    removeOption.onmouseover = () => removeOption.style.backgroundColor = '#f0f0f0';
    removeOption.onmouseout = () => removeOption.style.backgroundColor = 'transparent';
    removeOption.onclick = async () => {
        try {
            const response = await fetch('/api/friend-request', {
                method: 'PATCH',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
                body: JSON.stringify({
                    friendship_id: friendshipId,
                    status: 'rejected'
                })
            });
            
            if (response.ok) {
                alert(`Removed ${username} as friend`);
                dropdown.remove();
                loadFriends(); // Reload the list
            } else {
                alert('Failed to remove friend');
            }
        } catch (error) {
            console.error('Error removing friend:', error);
            alert('Error removing friend');
        }
    };
    
    dropdown.appendChild(removeOption);
    document.body.appendChild(dropdown);
    
    // Close dropdown when clicking outside
    const closeDropdown = function(e) {
        if (!dropdown.contains(e.target) && e.target !== nameSpan) {
            dropdown.remove();
            document.removeEventListener('click', closeDropdown);
        }
    };
    setTimeout(() => {
        document.addEventListener('click', closeDropdown);
    }, 10);
}

window.loadRooms = loadRooms;
window.loadFriends = loadFriends;