'use server'

import { cookies } from "next/headers";


export async function POST(req: Request, resp: Response) {

    // Get the URL from the environment variables
    const url = process.env.NEXT_PUBLIC_API_AUTH_LOGIN_URL;

    if (!url) {
        throw new Error('NEXT_PUBLIC_API_AUTH_LOGIN_URL is not defined');
    }

    const body = await req.json();

    const { email, password, lang, firebase_app_check_token, user  } = body;

    const res = await fetch(url, {
        method: 'POST',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            'X-Firebase-AppCheck': firebase_app_check_token,
        },
        body: JSON.stringify({
            email,
            password,
            lang,
            user
        }),
    })
    const data = await res.json()

    console.log("data XXX 111", data);

    // If the response is not successful, throw an error
    if (data.status == 'fail') {
        // TODO  Add phrase for error message in backend
       //  throw new Error(data.message);
    } 
    else {
        // Save access token in cookie
        cookies().set({
            name: 'access_token',
            value: data.access_token,
            httpOnly: true,
            path: '/',
            maxAge: 30 * 24 * 60 * 60, // 30 days
        })
    }

    // // Return the response
    console.log("data", data);

    return Response.json(data)

}
