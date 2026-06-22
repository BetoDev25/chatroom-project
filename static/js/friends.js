// static/js/friends.js
async function loadRooms() {
    console.log('loadRooms called');
    
    try {
        const response = await fetch('/api/rooms', {
            credentials: 'include'
        });

        console.log('Response status:', response.status);

        if (!response.ok) {
            console.error('Failed to load rooms');
            return;
        }

        const rooms = await response.json();
        console.log('Rooms received:', rooms);
        
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

// Auto-execute when script loads
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', loadRooms);
} else {
    loadRooms();
}