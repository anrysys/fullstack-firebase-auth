'use server'

import { cookies } from "next/headers";

export async function POST(req: Request, resp: Response) {

    // Get the URL from the environment variables
    const url = process.env.NEXT_PUBLIC_API_AUTH_USER_BY_EMAIL_URL;

    if (!url) {
        throw new Error('NEXT_PUBLIC_API_AUTH_USER_BY_EMAIL_URL is not defined');
    }

    const token = req.headers.get('X-Firebase-AppCheck');

    // Check if the token is defined
    if (!token) {
        throw new Error('X-Firebase-AppCheck is not defined');
    }

    // Get the email from the body of the request
    const body = await req.json();
    // console.log ('body', body);

    // Send the request to server to logout the user
    const response = await fetch(url, {
        method: 'POST',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            'X-Firebase-AppCheck': token,
            // add Bearer token to the header
            'Authorization': `Bearer ${cookies().get('access_token')?.value}`,
        },
        credentials: 'include', // attach cookies
        body: JSON.stringify(body),
    })
    
    return response;
}
