import { yupResolver } from '@hookform/resolvers/yup';
import { useForm } from 'react-hook-form';
import * as Yup from 'yup';

const profileSchema = Yup.object({
    name: Yup.string().nullable(),
    photo: Yup.string().nullable()
});

export const useProfileValidation = () => {
    const methods = useForm({
        resolver: yupResolver(profileSchema)
    });
    return methods;
};

const profilePasswordSchema = Yup.object({
    password: Yup.string().nullable(),   
});

export const useProfilePasswordValidation = () => {
    const methods = useForm({
        resolver: yupResolver(profilePasswordSchema)
    });
    return methods;
};