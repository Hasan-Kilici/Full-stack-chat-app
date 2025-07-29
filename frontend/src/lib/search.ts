export async function searchUsers(param: string, token: string){
    const res = await fetch(`http://localhost:4000/api/search/user?q=${param}`, {
    method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        }
    })
    
    const data = await res.json();

    return data
}