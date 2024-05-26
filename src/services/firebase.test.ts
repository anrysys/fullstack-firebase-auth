import { initializeApp } from "firebase/app";
import { ReCaptchaV3Provider } from "firebase/app-check";
import { getAuth } from "firebase/auth";
import { app } from "../services/firebase";

jest.mock("firebase/app", () => ({
  initializeApp: jest.fn(),
}));

jest.mock("firebase/auth", () => ({
  getAuth: jest.fn(),
}));

jest.mock("firebase/app-check", () => ({
  ReCaptchaV3Provider: jest.fn(),
}));

describe("Firebase Configuration", () => {
  test("should initialize the Firebase app with the correct configuration", () => {
    expect(initializeApp).toHaveBeenCalledWith({
      apiKey: process.env.NEXT_PUBLIC_FIREBASE_API_KEY,
      authDomain: process.env.NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN,
      databaseURL: process.env.NEXT_PUBLIC_FIREBASE_DATABASE_URL,
      projectId: process.env.NEXT_PUBLIC_FIREBASE_PROJECT_ID,
      storageBucket: process.env.NEXT_PUBLIC_FIREBASE_BUCKET,
      messagingSenderId: process.env.NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID,
      appId: process.env.NEXT_PUBLIC_FIREBASE_APP_ID,
      measurementId: process.env.NEXT_PUBLIC_FIREBASE_MEASUREMENT_ID,
    });
  });

  test("should get a reference to the auth service", () => {
    expect(getAuth).toHaveBeenCalledWith(app);
  });

  test("should initialize the reCAPTCHA v3 provider with the correct site key", () => {
    expect(ReCaptchaV3Provider).toHaveBeenCalledWith(
      process.env.NEXT_PUBLIC_RECAPTCHA_SITE_KEY
    );
  });
});


// describe("firebaseConfig", () => {
// const firebaseConfig = {
//     apiKey: process.env.NEXT_PUBLIC_FIREBASE_API_KEY,
//     authDomain: process.env.NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN,
//     databaseURL: process.env.NEXT_PUBLIC_FIREBASE_DATABASE_URL,
//     projectId: process.env.NEXT_PUBLIC_FIREBASE_PROJECT_ID,
//     storageBucket: process.env.NEXT_PUBLIC_FIREBASE_BUCKET,
//     messagingSenderId: process.env.NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID,
//     appId: process.env.NEXT_PUBLIC_FIREBASE_APP_ID,
//     measurementId: process.env.NEXT_PUBLIC_FIREBASE_MEASUREMENT_ID,
// };

// it("should have all the required properties", () => {
//     expect(firebaseConfig).toHaveProperty("apiKey");
//     expect(firebaseConfig).toHaveProperty("authDomain");
//     expect(firebaseConfig).toHaveProperty("databaseURL");
//     expect(firebaseConfig).toHaveProperty("projectId");
//     expect(firebaseConfig).toHaveProperty("storageBucket");
//     expect(firebaseConfig).toHaveProperty("messagingSenderId");
//     expect(firebaseConfig).toHaveProperty("appId");
//     expect(firebaseConfig).toHaveProperty("measurementId");
// });

//   it("should have valid values for the properties", () => {
//     expect(typeof firebaseConfig.apiKey).toBe("undefined");
//     expect(typeof firebaseConfig.authDomain).toBe("string");
//     expect(typeof firebaseConfig.databaseURL).toBe("string");
//     expect(typeof firebaseConfig.projectId).toBe("string");
//     expect(typeof firebaseConfig.storageBucket).toBe("string");
//     expect(typeof firebaseConfig.messagingSenderId).toBe("string");
//     expect(typeof firebaseConfig.appId).toBe("string");
//     expect(typeof firebaseConfig.measurementId).toBe("string");
//   });
// });