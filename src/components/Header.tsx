"use client";
import { HOME_ROUTE, LOGIN_ROUTE, PROFILE_ROUTE, REGISTER_ROUTE } from "@/constants/routes";
import { AuthContext } from "@/provider/AuthProvider";
import { auth } from "@/services/firebase";
import { signOut } from "firebase/auth";
import Link from "next/link";
import { useRouter } from "next/navigation";

const Header = () => {
    const { user }: any = AuthContext();
    const router = useRouter();
    const logOut = () => {
        signOut(auth).then(async (objResponseFromFirebase) => {

        // // Fetch locale from user browser & set it to the values object
        // const values: any = {}; // Declare the 'values' variable
        // const locale = navigator.language; // "en-US"
        // const lang = locale.slice(0, 2); // "en"
        // values.lang = lang;
       
         console.log("Logout: objResponseFromFirebase", objResponseFromFirebase);
        // values.user = objResponseFromFirebase.user;
        // values.user.isLogin = false;
        // values.user.email = "";
        // values.user.password = "";

        // console.log("values + objResponseFromFirebase ", values);
        // Save the authentication result to the RESTful API
        const response = await fetch('/api/auth/logout', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            //body: JSON.stringify(objResponseFromFirebase),
        });
            //console.log("response FROM API", response);

            // if (response.ok) {
            //     alert("User Logout Successfully");

            //     router.push(PROFILE_ROUTE);
            // } else {
            //     console.log("catch ", response.statusText);
            //     alert("Something went wrong please try again");
            //     //throw new Error('Failed to save authentication result to the RESTful API');
            // }

            router.push(LOGIN_ROUTE);
        }).catch((e) => {
            console.log("Logout Catch ", e.message)
        })
    }

    return (
        <header className="h-20 bg-gradient-to-br from-yellow-400/20 via-blue-300 to-purple-400/60 flex px-10 drop-shadow-[0px_2px_10px_rgba(2,0,0) text-black">
            <nav className="w-full mx-auto flex justify-between items-center px-2 text-black font-serif text-xl">
                <Link href={HOME_ROUTE}><div>Logo</div></Link>
                <ul className="flex gap-4">
                    {!user?.isLogin &&
                        <>
                            <Link href={LOGIN_ROUTE}><li>Login</li></Link>
                            <Link href={REGISTER_ROUTE}><li>Register</li></Link>
                        </>
                    }
                    {user?.isLogin &&
                        <>
                            <Link href={PROFILE_ROUTE}><li>Profile</li></Link>
                            <li className=" cursor-pointer" onClick={logOut}>Logout</li>
                        </>
                    }
                </ul>
            </nav>
        </header>
    )
}

export default Header;