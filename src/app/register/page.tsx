"use client";
import SubmitButton from "@/components/Button";
import InputField from "@/components/InputField";
import { LOGIN_ROUTE, PROFILE_ROUTE } from "@/constants/routes";
import useAuthentication from "@/hooks/useAuthentication";
import { app, auth, provider } from '@/services/firebase';
// import { auth } from '@/services/firebase';
import { useRegisterValidation } from "@/validationSchema/useAuth";
import { createUserWithEmailAndPassword } from "firebase/auth";
import Link from "next/link";
import { useRouter } from "next/navigation";

import { getToken, initializeAppCheck } from 'firebase/app-check';

const Register = () => {
    const router = useRouter();
    useAuthentication();
    const { handleSubmit, register, formState: { errors }, reset } = useRegisterValidation();

    // Initialize the Firebase App Check
   const appCheck = initializeAppCheck(app, { provider: provider });

    const submitForm = async (values: any) => {

        // Fetch locale from user browser & set it to the values object
        const locale = navigator.language; // "en-US"
        const lang = locale.slice(0, 2); // "en"
        values.lang = lang;
        const { email, password, cnfPassword } = values;

        createUserWithEmailAndPassword(auth, values.email, values.password).then( async (objResponseFromFirebase) => {

            let appCheckTokenResponse;
            try {
                appCheckTokenResponse = await getToken(appCheck, /* forceRefresh= */ false);
            } catch (err) {
                // Handle any errors if the token was not retrieved.
                console.error(err);
                return;
            }

            // Add objResponseFromFirebase to values
            values.user = objResponseFromFirebase.user;
            // Add the Firebase App Check token to the values object
            values.firebase_app_check_token = appCheckTokenResponse.token

            // Save the authentication result to the RESTful API
            const response = await fetch('/api/auth/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(values),
            });

            // Convert the response to JSON
            const resp = await response.json();

            console.log("Response data from API", resp);

            // Handle the response
            if (resp.status == 'success') {
                alert("User Register Successfully");
                reset();
                router.push(PROFILE_ROUTE);
            } else {
                // console.log("catch ", response.statusText);
                alert(resp.errors);
                //throw new Error('Failed to save authentication result to the RESTful API');
                router.push(LOGIN_ROUTE);
            }

        }).catch(e => {
            console.log("catch ", e.message);
            alert("Something went wrong please try again");
        })
    }

    return (
        <div className="h-screen flex justify-center items-center bg-gradient-to-br from-yellow-400/20 via-blue-300 to-purple-400/60">
            <div className="w-1/2 rounded-md bg-white/30 shadow-lg flex justify-between flex-col">
                <div className="h-28 w-full justify-center flex items-center">
                    <span className="text-3xl text-black font-mono font-semibold bg-yellow-300 p-3 rounded-lg">Welcome To Register</span>
                </div>
                <form onSubmit={handleSubmit(submitForm)} className="h-full w-1/2 mx-auto ">
                    <InputField
                        register={register}
                        error={errors.email}
                        type="text"
                        placeholder="Enter Your Email Here..."
                        name="email"
                        label="Email"
                    />
                    <InputField
                        register={register}
                        error={errors.password}
                        type="password"
                        placeholder="Enter Your Password Here..."
                        name="password"
                        label="Password"
                    />
                    <InputField
                        register={register}
                        error={errors.cnfPassword}
                        type="password"
                        placeholder="Enter Your Confirm Password Here..."
                        name="cnfPassword"
                        label="Confirm Password"
                    />
                    <SubmitButton label="Submit" />
                </form>
                <div className="h-20 mx-auto">
                    <span className="text-sm text-gray-600">Already have account?
                        <Link href={LOGIN_ROUTE}><span className="text-blue-500 font-semibold text-md" > Login Here</span></Link>
                    </span>
                </div>
            </div>
        </div>
    )
}

export default Register;