'use server'

export async function POST(req: Request, resp: Response) {

    const url = process.env.NEXT_PUBLIC_API_AUTH_REGISTER_URL;

    if (!url) {
        throw new Error('NEXT_PUBLIC_API_AUTH_REGISTER_URL is not defined');
    }

    // Delete cookie access token on api register
    const res = await fetch(url, {
        method: 'GET',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
        },
        // body: JSON.stringify({
        //     lang
        // }),        
    })
    const data = await res.json()

    // If the response is not successful, throw an error
    if (data.status !== 'success') {
        // TODO  Add phrase for error message in backend
        // throw new Error(data.message);
    } else {
        // Delete cookie access token
        // cookies().delete('access_token')

    }

    return Response.json(data)

}
