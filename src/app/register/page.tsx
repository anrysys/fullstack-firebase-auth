"use client";

import SubmitButton from "@/components/Button";
import InputField from "@/components/InputField";
import { FORGOT_PASSWORD_ROUTE, LOGIN_ROUTE, PROFILE_ROUTE } from "@/constants/routes";
import useAuthentication from "@/hooks/useAuthentication";
import { app, auth, providers } from '@/services/firebase';
import { useRegisterValidation } from "@/validationSchema/useAuth";
import { getToken, initializeAppCheck } from 'firebase/app-check';
import { createUserWithEmailAndPassword, getAdditionalUserInfo } from "firebase/auth";
import Link from "next/link";
import { useRouter } from "next/navigation";

interface UserData {
    uid: string;
    providerId: string;
    tenantId: string;
    isNewUser: boolean;
    isAnonymous: boolean;
    displayName: string;
    email: string;
    emailVerified: boolean;
    phoneNumber: string;
    photoURL: string;
    lang: string;
    providerData?: any;
}

const Register = () => {
    const router = useRouter();
    useAuthentication();
    const { handleSubmit, register, formState: { errors }, reset } = useRegisterValidation();

    // Initialize the Firebase App Check
    const appCheck = initializeAppCheck(app, { provider: providers.reCaptchProvider });

    const submitForm = async (values: any) => {

        // Fetch locale from user browser & set it to the values object
        const locale = navigator.language; // "en-US"
        const lang = locale.slice(0, 2); // "en"
        values.lang = lang;
        auth.languageCode = lang;
        const { email, password, cnfPassword } = values;

        createUserWithEmailAndPassword(auth, values.email, values.password).then(async (result) => {

            let firebase_app_check_token;
            try {
                firebase_app_check_token = await getToken(appCheck, /* forceRefresh= */ false);
            } catch (err) {
                // Handle any errors if the token was not retrieved.
                console.error(err);
                return;
            }

            // The signed-in user info.
            const user = result.user;

            // IdP data available using getAdditionalUserInfo(result).profile
            const additionalUserInfo = getAdditionalUserInfo(result);
            console.log("GitHub login successful:", user);
            console.log("GitHub additionalUserInfo:", additionalUserInfo);

            // Get the user data
            const userData: UserData = {
                uid: user.uid,
                providerId: 'github.com',
                tenantId: user.tenantId as string,
                isNewUser: additionalUserInfo?.isNewUser as boolean,
                isAnonymous: user.isAnonymous,
                displayName: user.displayName as string,
                email: user.email as string,
                emailVerified: user.emailVerified,
                phoneNumber: user.phoneNumber as string,
                photoURL: user.photoURL as string,
                lang: navigator.language.slice(0, 2) as string,
                providerData: user.providerData,
            }
            userData.providerId = additionalUserInfo?.providerId as string; // 'github.com'
            userData.isNewUser = additionalUserInfo?.isNewUser as boolean;

            console.log("NEXT USER DATA: ", userData);

            // Save the authentication result to the RESTful API
            const response = await fetch('/api/auth/login-socials', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'X-Firebase-AppCheck': firebase_app_check_token.token,
                },
                body: JSON.stringify(
                    userData
                ),
            });

            // Convert the response to JSON
            const resp = await response.json();

            console.log("response on client: ", resp);

            // Handle the response
            if (resp.status == 'success') {
               // alert("User Register Successfully");
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
                    <span className="text-3xl text-black font-mono font-semibold bg-yellow-300 p-3 rounded-lg">
                        Register
                    </span>
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
                    <p className="text-sm text-gray-600">
                        <Link href={FORGOT_PASSWORD_ROUTE}><span className="text-blue-500 font-semibold text-md">Forgot Password?&nbsp;</span></Link>
                    </p>
                </div>
            </div>
        </div>
    )
}

export default Register;