// UserDataRepo.ts

import { app } from '@/services/firebase';
import { collection, doc, getDoc, getFirestore, setDoc } from 'firebase/firestore';

const storage = getFirestore(app);

interface UserDataRepoInterface {
    get(user: any): Promise<any>;
    set(user: any, data: any): void;
    delete(user: any): void;
    merge(prevUserData: any, currentUserData: any): any;
}

class UserDataRepo implements UserDataRepoInterface {
    
    async get(user: any) {
        if (user) {
            const userDocRef = doc(collection(storage, 'users'), user.uid);
            const userDocSnap = await getDoc(userDocRef);
            if (userDocSnap.exists()) {
                return userDocSnap.data();
            }
        }
        return null;
    }

    async set(user: any, data: any) {
        if (user) {
            const userDocRef = doc(collection(storage, 'users'), user.uid);

            // Convert the custom UserMetadata object into a plain object
            // Assuming `data` contains the UserMetadata object
            // This is a simple example; you might need a more complex approach depending on the structure of your data
            const plainData = this.convertToPlainObject(data);

            await setDoc(userDocRef, plainData);
        }
    }


    delete(user: any) {
        if (user) {
            const userDocRef = doc(collection(storage, 'users'), user.uid);
            // Implement the logic to delete user data from your storage
            console.log("Deleting user data for user:", user.uid);
        }
    }

    merge(prevUserData: any, currentUserData: any) {
        const mergedData = {
            ...prevUserData,
            ...currentUserData,
        };
        return mergedData;
    }
    // Helper method to convert a custom object into a plain object
    convertToPlainObject(obj: any): any {
        // Implement a more sophisticated conversion if needed
        // This is a simple approach and may not work for deeply nested or complex objects
        return JSON.parse(JSON.stringify(obj));
    }

}

export default UserDataRepo;