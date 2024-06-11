import { initializeApp } from "firebase/app";
import { ReCaptchaV3Provider } from "firebase/app-check";
import { EmailAuthProvider, FacebookAuthProvider, GithubAuthProvider, GoogleAuthProvider, TwitterAuthProvider, getAuth } from "firebase/auth";

// !!! IMPORTANT !!!
// !Add secret key (in field 'SECRET KEY': ************ ) to project settings in Firebase Console https://console.firebase.google.com/project/YOUR-PROJECT/appcheck/apps

// TODO: Add SDKs for Firebase products that you want to use
// https://firebase.google.com/docs/web/setup#available-libraries

// Your web app's Firebase configuration
// Attention Please add your firebase console config here
const firebaseConfig = {
  apiKey: process.env.NEXT_PUBLIC_FIREBASE_API_KEY,
  authDomain: process.env.NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN,
  databaseURL: process.env.NEXT_PUBLIC_FIREBASE_DATABASE_URL,
  projectId: process.env.NEXT_PUBLIC_FIREBASE_PROJECT_ID,
  storageBucket: process.env.NEXT_PUBLIC_FIREBASE_BUCKET,
  messagingSenderId: process.env.NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID,
  appId: process.env.NEXT_PUBLIC_FIREBASE_APP_ID,
  measurementId: process.env.NEXT_PUBLIC_FIREBASE_MEASUREMENT_ID,
};

// Initialize Firebase
const app = initializeApp(firebaseConfig);

// Get a reference to the auth service
const auth = getAuth(app);

// Initialize providers
const providers = {
  emailAuthProvider: new EmailAuthProvider(),
  reCaptchProvider: new ReCaptchaV3Provider(process.env.NEXT_PUBLIC_RECAPTCHA_SITE_KEY as string),
  googleProvider: new GoogleAuthProvider(),
  facebookProvider: new FacebookAuthProvider(),
  twitterProvider: new TwitterAuthProvider(),
  githubProvider: new GithubAuthProvider(),
};

// Get a reference to the Firestore service
// const db = getFirestore(app);

// Get a reference to the storage service, which is used to create references in your storage bucket
export { app, auth, providers };

