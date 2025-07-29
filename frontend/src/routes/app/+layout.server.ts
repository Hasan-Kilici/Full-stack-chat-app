import { redirect } from '@sveltejs/kit';
import { getAuth, getUserRooms } from '$lib/utils';
export async function load(event: any) {
  const { cookies } = event;
  const token = cookies.get('jwt_token');

  if (!token) {
    throw redirect(302, '/login');
  }


  const userData = await getAuth(token);
  const roomsData = await getUserRooms(token);

  return {
    user: {
        name: userData.name,
        email: userData.email,
        user_id: userData.user_id,
        auth: token,
     },
     rooms: roomsData
  };
    
}