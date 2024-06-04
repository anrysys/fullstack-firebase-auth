'use server'


// Path: src/app/api/auth/reset-password/route.ts
// https://YOUR-PROJECT.firebaseapp.com/__/auth/action?mode=resetPassword&oobCode=03nw-L6pMG48hK7B8MAzah2cQ5r6jAXm_oE3xptHURUAAAGP3SYpUQ&apiKey=AIzaSyAjInYaz6vcxH01PwqOu9fLHSUjOa8eDA8&lang=en
export async function GET(req: Request, resp: Response) {

    // Get the query string from the request URL
    const url = new URL(req.url);
    const searchParams = new URLSearchParams(url.search);


    console.log("searchParams", searchParams);



    // Get the values from the query string
    const mode = searchParams.get('mode');
    const oobCode = searchParams.get('oobCode');

    // Use the values as needed
    console.log("mode", mode);
    console.log("oobCode", oobCode);

    // Check if the mode is for resetting the password
    if (mode === 'resetPassword') {
        // Redirect to the password reset page
        return Response.redirect(`http://localhost:3000/reset-password?oobCode=${oobCode}`);
    }

    // Redirect to the home page
    return Response.redirect(`http://localhost:3000/`);

    // return new Response('oobCode: ' + oobCode + '\n\n mode: ' + mode + '\n\n', { status: 200 })

}
