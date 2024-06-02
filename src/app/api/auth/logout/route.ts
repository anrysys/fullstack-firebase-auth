'use server'

import { cookies } from "next/headers";

export async function GET(req: Request, resp: Response) {

    // Get the URL from the environment variables
    const url = process.env.NEXT_PUBLIC_API_AUTH_LOGOUT_URL;

    if (!url) {
        throw new Error('NEXT_PUBLIC_API_AUTH_LOGOUT_URL is not defined');
    }

    const token = req.headers.get('X-Firebase-AppCheck');

    // Check if the token is defined
    if (!token) {
        throw new Error('X-Firebase-AppCheck is not defined');
    }

    // Send the request to server to logout the user
    const response = await fetch(url, {
        method: 'GET',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            'X-Firebase-AppCheck': token,
            // add Bearer token to the header
            'Authorization': `Bearer ${cookies().get('access_token')?.value}`,
        },
        credentials: 'include', // отправить куки
    })
        .then(response => response.json())
        //.then(data => console.log(data))
        .catch((error) => {
            // console.error('Error Logout:', error);
        });


    // Delete access token in cookie
    cookies().set({
        name: 'access_token',
        value: '',
        httpOnly: true,
        path: '/',
        maxAge: 0, // 30 days
    })
    // delete refresh token in cookie
    cookies().set({
        name: 'refresh_token',
        value: '',
        httpOnly: true,
        path: '/',
        maxAge: 0, // 30 days
    })

    return new Response('Logout', { status: 200 })
}
