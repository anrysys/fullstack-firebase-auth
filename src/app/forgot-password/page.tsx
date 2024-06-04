"use client";
import SubmitButton from "@/components/Button";
import InputField from "@/components/InputField";
import { FORGOT_PASSWORD_ROUTE, REGISTER_ROUTE } from "@/constants/routes";
import useAuthentication from "@/hooks/useAuthentication";
import { app, auth, provider } from '@/services/firebase';
import { useForgotPasswordValidation } from "@/validationSchema/useAuth";
import { getToken, initializeAppCheck } from 'firebase/app-check';
import { sendPasswordResetEmail } from "firebase/auth";
import Link from "next/link";
import { useRouter } from "next/navigation";

const ForgotPassword = () => {
    const { handleSubmit, register, formState: { errors }, reset } = useForgotPasswordValidation();
    const router = useRouter();
    useAuthentication();

    // Initialize the Firebase App Check
    const appCheck = initializeAppCheck(app, { provider: provider });

    const submitForm = async (values: any) => {

        // Fetch locale from user browser & set it to the values object
        const locale = navigator.language; // "en-US"
        const lang = locale.slice(0, 2); // "en"
        values.lang = lang;
        auth.languageCode = lang;
        const { email } = values;

        // Send password reset email
        // Если вам нужно знать, когда пользователь сбросил свой пароль, 
        // вы можете использовать функцию 'onAuthStateChanged' из Firebase, 
        // чтобы слушать изменения в состоянии аутентификации пользователя.
        await sendPasswordResetEmail(auth, values.email, {
            // !!! TODO: Change to your domain
            url: 'http://localhost:3000/login', // url: 'http://localhost/api/auth/action' (API_AUTH_ACTION_ROUTE), 
            handleCodeInApp: true
        }).then(async (objResponseFromFirebase) => {

            console.log("Password reset email sent successfully", objResponseFromFirebase);

            // confirm Password Reset
            let appCheckTokenResponse;
            try {
                appCheckTokenResponse = await getToken(appCheck, /* forceRefresh= */ false);
            } catch (err) {
                // Handle any errors if the token was not retrieved.
                console.error('Error getting token:', err);
            }

            alert("Password reset email sent successfully");
            reset();
            // router.push('/login');
        })
            .catch((error) => {
                // const errorCode = error.code;
                // const errorMessage = error.message;
                // ..
            });

    }
    // TODO: DELETE THESE DEFAULT VALUES !!!
    const emailDefaultValue = "anrysys@gmail.com";

    return (
        <div className="h-screen flex justify-center items-center bg-gradient-to-br from-yellow-400/20 via-blue-300 to-purple-400/60">
            <div className="w-1/2 rounded-md bg-white/30 shadow-lg flex justify-between flex-col">
                <div className="h-28 w-full justify-center flex items-center">
                    <span className="text-3xl text-black font-mono font-semibold bg-yellow-300 p-3 rounded-lg">Forgot Password</span>
                </div>
                <form onSubmit={handleSubmit(submitForm)} className="h-full w-1/2 mx-auto ">
                    <InputField
                        register={register}
                        error={errors.email}
                        type="text"
                        placeholder="Enter Your Email Here..."
                        name="email"
                        label="Email"
                        defaultValue={emailDefaultValue}

                    />
                    <SubmitButton label="Submit" />
                </form>
                <div className="h-20 mx-auto">
                    <span className="text-sm text-gray-600">Dont have an account?
                        <Link href={REGISTER_ROUTE}><span className="text-blue-500 font-semibold text-md">Register Here</span></Link>
                    </span>
                    <p className="text-sm text-gray-600">
                        <Link href={FORGOT_PASSWORD_ROUTE}><span className="text-blue-500 font-semibold text-md">Forgot Password?&nbsp;</span></Link>
                    </p>
                </div>
            </div>
        </div>
    )
}

export default ForgotPassword;
// Path: src/app/reset-password/page.tsx