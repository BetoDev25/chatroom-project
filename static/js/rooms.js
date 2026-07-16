function debounce(func, wait) {
    let timeout;
    return function(...args) {
        clearTimeout(timeout);
        timeout = setTimeout(() => func.apply(this, args), wait);
    };
}

async function search(query) {
    const rooms = await callRooms(query || '');
    loadPublicRooms(rooms);
}
async function callRooms(roomName) {
    try {
        const url = roomName ? `/api/rooms-public/${roomName}` : '/api/rooms-public/';
        const response = await fetch(url, {
            credentials: 'include'
        });

        if (!response.ok) {
            console.error('Failed to load public rooms');
            return
        }

        let rooms = await response.json();
        return rooms;
    } catch (error) {
        console.error('Error loading rooms:', error);
        return [];
    }
}

function loadPublicRooms(rooms) {
    const roomList = document.getElementById('publicRoomList');
    if (!roomList) {
        console.error('publicRoomList element not found');
        return;
    }
    
    roomList.innerHTML = '';
    if (!rooms) {
        roomList.innerHTML = '<div style="color: #888; padding: 20px 0; text-align: center;">Sorry, that room could not be found.</div>';
        return;
    }
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

        if (currentUser.id === room.OwnerID) {
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
            buttonContainer.appendChild(deleteButton);
        }

        buttonContainer.appendChild(joinButton);

        div.appendChild(nameSpan);
        div.appendChild(buttonContainer);
        roomList.appendChild(div);
    });
}

function showPublicRooms() {
    currentRoom = null;
    currentConversation = null;
    const messagesDiv = document.getElementById('messages');
    fetch('/static/rooms.html')
        .then(res => res.text())
        .then(html => {
            messagesDiv.innerHTML = html;
            setupSearch();
            search('');
        })
        .catch(err => console.error('Error loading rooms panel:', err));
}

function setupSearch() {
    const searchInput = document.getElementById('roomSearch');
    if (!searchInput) {
        console.error('Search input not found');
        return;
    }
    
    const debounceSearch = debounce(async () => {
        const query = searchInput.value.trim();
        await search(query);
    }, 300);
    
    searchInput.addEventListener('input', debounceSearch);
}
