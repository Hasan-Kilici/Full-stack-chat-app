import { redirect } from "@sveltejs/kit";

export async function getAuth(token: string) {
    const res = await fetch('http://localhost:4000/auth/@me', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        }
    })

  const userData = await res.json();

  return userData;
}

export async function getUserRooms(token: string) {
    const res = await fetch('http://localhost:4000/api/get/participants/', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        }
    })

    const roomsData = await res.json();

    return roomsData;
}

export async function getMessages(token: string, room_id: string) {
    const res = await fetch(`http://localhost:4000/api/get/messages/${room_id}`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        }
    })

    const messagesData = await res.json();

    return messagesData;
}

export async function createRoom (token: string, other_user_id: string) {
    const res = await fetch(`http://localhost:4000/room/room`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
            is_group:false,
            other_user_id,
            name:"",
        })
    })

    const roomData = await res.json()
    
    window.location.href =  `/app/chat/${roomData.room_id}`
}