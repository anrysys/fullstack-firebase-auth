"use client";
import SubmitButton from "@/components/Button";
import ButtonProps from "@/components/ButtonProps";
import InputField from "@/components/InputField";
import { FORGOT_PASSWORD_ROUTE, PROFILE_ROUTE, REGISTER_ROUTE } from "@/constants/routes";
import useAuthentication from "@/hooks/useAuthentication";
import { app, auth, providers } from '@/services/firebase';
import { useLoginValidation } from "@/validationSchema/useAuth";
import { getToken, initializeAppCheck } from 'firebase/app-check';
import { GithubAuthProvider, User, fetchSignInMethodsForEmail, getAdditionalUserInfo, signInWithEmailAndPassword, signInWithPopup } from "firebase/auth";
import Link from "next/link";
import { useRouter } from "next/navigation";
import router from "next/router";

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

const handleGitHubLogin = async () => {

    // Initialize the Firebase App Check
    const appCheck = initializeAppCheck(app, { provider: providers.reCaptchProvider });
    const firebase_app_check_token = await getToken(appCheck);

    try {
        // Sign in with GitHub
        const result = await signInWithPopup(auth, providers.githubProvider).then(async (result) => {

            // This gives you a Google Access Token. You can use it to access the Google API.
            const credential = GithubAuthProvider.credentialFromResult(result);
            const token = credential?.accessToken;
            // The signed-in user info.
            const user = result.user;
            // IdP data available using getAdditionalUserInfo(result).profile
            const additionalUserInfo = getAdditionalUserInfo(result);
            console.log("GitHub login successful:", user);
            console.log("GitHub token:", token);
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

            // Check by email if the user is new or existing and redirect accordingly to the profile page or the login page
            // Check for the presence of a user in the backend database by email
            const response = await fetch('/api/auth/login-socials', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'X-Firebase-AppCheck': firebase_app_check_token.token,
                },
                //credentials: 'include', // отправить куки
                body: JSON.stringify(
                    userData,
                )
            })

            const resp = await response.json();

            console.log("response on client: ", resp);

            if (resp.status == 'success') {
                // alert("Hello, " + userData.displayName + "!");
                // Redirect to the profile page
                router.push(PROFILE_ROUTE);
            }

            // ...
        }).catch((error) => {
            // Handle Errors here.
            const errorCode = error.code;
            const errorMessage = error.message;
            // The email of the user's account used.
            const email = error.email;
            // The AuthCredential type that was used.
            const credential = GithubAuthProvider.credentialFromError(error);

            // If the error code is 'auth/missing-identifier', the user is not signed in. 
            if (errorCode === 'auth/account-exists-with-different-credential') {
                alert('Вы пытаетесь войти с помощью другого аккаунта. Пожалуйста, войдите с помощью того же аккаунта, который вы использовали при первоначальной регистрации.')
                console.log("The user is not signed in");
            }
        });
    } catch (error: any) {

        console.error("Error signing in with GitHub 1:", error);

        const user = auth.currentUser as User;
        // Получаем email, связанный с GitHub аккаунтом
        const email = error.customData.email;
        if (error.code === 'auth/account-exists-with-different-credential') {

            try {
                // Получаем список методов аутентификации, доступных для этого email
                const signInMethods = await fetchSignInMethodsForEmail(auth, email);
                // Проверяем доступные методы и предлагаем пользователю войти
                // с использованием одного из них
                let msg = "Пожалуйста, войдите, используя ваш GitHub аккаунт";
                if (signInMethods.includes('password')) {
                    // TODO: В профиле добавить возможность привязки GitHub аккаунта (в настройках профиля в разделе 'Связанные аккаунты')
                    alert("Пожалуйста, войдите, используя ваш email и пароль и затем добавьте Ваш GitHub аккаунт в настройках профиля в разделе 'Связанные аккаунты'")
                } else if (signInMethods.includes('github.com')) {
                    console.log("Пожалуйста, войдите, используя ваш GitHub аккаунт");
                }
                alert(msg)
                // Добавьте дополнительные условия для других методов аутентификации, если необходимо
            } catch (fetchError) {
                console.error("Ошибка при получении методов аутентификации:", fetchError);
            }
        }
        else {
            console.error("Error signing in with GitHub 2:", error);
        }
    }

};


const Login = () => {
    const { handleSubmit, register, formState: { errors }, reset } = useLoginValidation();
    const router = useRouter();
    useAuthentication();

    // Fetch locale from user browser & set it to the values object
    const locale = navigator.language; // "en-US"
    const lang = locale.slice(0, 2); // "en"
    // auth.languageCode = lang;

    const submitForm = async (values: any) => {

        values.lang = lang;
        const { email, password } = values;

        _signInWithEmailAndPassword(values, reset, router);

    }
    // TODO: DELETE THESE DEFAULT VALUES !!!
    const emailDefaultValue = "anrysys@gmail.com";
    const passwordDefaultValue = "anrysys@gmail.com";

    return (
        <div className="h-screen flex justify-center items-center bg-gradient-to-br from-yellow-400/20 via-blue-300 to-purple-400/60">
            <div className="w-1/2 rounded-md bg-white/30 shadow-lg flex justify-between flex-col">
                <div className="h-28 w-full justify-center flex items-center">
                    <span className="text-3xl text-black font-mono font-semibold bg-yellow-300 p-3 rounded-lg">Login</span>
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
                    <InputField
                        register={register}
                        error={errors.password}
                        type="password"
                        placeholder="Enter Your Password Here..."
                        name="password"
                        label="Password"
                        defaultValue={passwordDefaultValue}
                    />
                    <SubmitButton label="Submit" />
                </form>

                <ButtonProps buttonText="Login with GitHub" backgroundColor="green-500" hoverBackgroundColor="green-600" textColor="white" onClick={handleGitHubLogin} />

                <div className="h-20 mx-auto">
                    <span className="text-sm text-gray-600">Dont have an account?&nbsp;
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
export default Login;

async function _signInWithEmailAndPassword(values: any, reset: any, router: any) {

    try {
        // Sign in with email and password

        signInWithEmailAndPassword(auth, values.email, values.password).then(async (objResponseFromFirebase) => {

            // Initialize the Firebase App Check
            const appCheck = initializeAppCheck(app, { provider: providers.reCaptchProvider });

            let firebase_app_check_token;
            try {
                firebase_app_check_token = await getToken(appCheck, /* forceRefresh= */ false);
            } catch (err) {
                // Handle any errors if the token was not retrieved.
                console.error(err);
                return;
            }

            // Add objResponseFromFirebase to values
            values.user = objResponseFromFirebase.user;
            // Add the Firebase App Check token to the values object
            values.firebase_app_check_token = firebase_app_check_token.token;

            // Save the authentication result to the RESTful API
            const response = await fetch('/api/auth/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(values),
            });

            // Convert the response to JSON
            const resp = await response.json();

            // Handle the response
            if (resp.status == 'success') {
                alert("User Login Successfully");
                reset();
                router.push(PROFILE_ROUTE);
            } else {
                // console.log("catch ", response.statusText);
                alert(resp.errors);
                //throw new Error('Failed to save authentication result to the RESTful API');
            }

        }).catch((e) => {
            console.log("Login Error ", e.message);
            const errorCode = e.code;
            // If the error code is 'auth/missing-identifier', the user is not signed in. 
            if (errorCode === 'auth/account-exists-with-different-credential'
                || errorCode === 'auth/wrong-password'
            ) {
                alert('Вы пытаетесь войти с помощью другого аккаунта. Пожалуйста, войдите с помощью того же аккаунта, который вы использовали при первоначальной регистрации.')
                console.log("The user is not signed in");
            }

        });

    } catch (error: any) {

        const user = auth.currentUser as User;
        // Получаем email, связанный с GitHub аккаунтом
        const email = error.customData.email;
        if (error.code === 'auth/account-exists-with-different-credential') {

            try {
                // Получаем список методов аутентификации, доступных для этого email
                const signInMethods = await fetchSignInMethodsForEmail(auth, email);
                // Проверяем доступные методы и предлагаем пользователю войти
                // с использованием одного из них
                let msg = "Пожалуйста, войдите, используя ваш GitHub аккаунт";
                if (signInMethods.includes('password')) {
                    // TODO: В профиле добавить возможность привязки GitHub аккаунта (в настройках профиля в разделе 'Связанные аккаунты')
                    alert("Пожалуйста, войдите, используя ваш email и пароль и затем добавьте Ваш GitHub аккаунт в настройках профиля в разделе 'Связанные аккаунты'")
                } else if (signInMethods.includes('github.com')) {
                    console.log("Пожалуйста, войдите, используя ваш GitHub аккаунт");
                }
                alert(msg)
                // Добавьте дополнительные условия для других методов аутентификации, если необходимо
            } catch (fetchError) {
                console.error("Ошибка при получении методов аутентификации:", fetchError);
            }
        }
        else {
            console.error("Error signing in with GitHub:", error);
        }
    }

}
