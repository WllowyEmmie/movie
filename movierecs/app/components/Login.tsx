'use client'
import React, { useState } from "react";
import axios from "axios";
import { useRouter } from "next/navigation";
interface LoginForm {
    email: string;
    password: string
}
type FormProps = {
    routes: string,
    destination: string,
}
function LoginForm({ routes, destination }: FormProps) {
    const [form, setForm] = useState<LoginForm>({
        email: "",
        password: "",
    });
    const router = useRouter()
    const [error, setError] = useState<string | null>(null)
    const validateForm = () => {

        if (!form.email?.trim() || !form.password.trim()) {
            return ("All fields are required");
        }
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (!emailRegex.test(form.email)) {
            return ("Input a valid email");
        }


        const passwordRegex = /^(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*]).{6,}$/;
        if (!passwordRegex.test(form.password)) {
            return ("Password Should have a Capital letter a number and a symbol and must be at least 6 characters");
        }

    };
    const handleSubmit = async (event: React.FormEvent) => {
        event.preventDefault();
        const validateError = validateForm();
        if (validateError) {
            setError(validateError)
            return
        }
        try {
            const response = await axios.post(routes, form);

            const data = response.data;
            const token = data.token;
            const userID = data.user_id
            localStorage.setItem("token", token);
            localStorage.setItem("userID", userID);

            router.push(destination as string)
        } catch (error: any) {
            setError(error.response?.data?.error || "Failed")
        }
    }
    return (
        <div className="flex items-center justify-center min-h-screen bg-darkPurple">

            <form
                onSubmit={handleSubmit}
                className="bg-blackish p-8 rounded-lg shadow-lg w-full max-w-md"
            >
                <h2 className="text-2xl font-bold mb-6 text-center text-lightPurple">
                    Login
                </h2>

                {error && (
                    <div className="text-center bg-red-100 text-600 text-error p-2 mb-4 rounded text-sm">
                        {error}
                    </div>
                )}

                <div className="mb-4">
                    <label
                        htmlFor="email"
                        className="block mb-1 text-sm font-medium text-white"
                    >
                        Email
                    </label>
                    <input
                        type="email"
                        id="email"
                        value={form.email}
                        onChange={(e) => setForm({ ...form, email: e.target.value })}
                        className="w-full px-3 py-2 focus:text-blackish rounded focus:outline-none focus:ring focus:bg-lightPurple focus:border-lightPurple focus: bg-darkPurple"
                        required
                    />
                </div>

                <div className="mb-4">
                    <label
                        htmlFor="password"
                        className="block mb-1 text-sm font-medium text-white"
                    >
                        Password
                    </label>
                    <input
                        type="password"
                        id="password"
                        value={form.password}
                        onChange={(e) => setForm({ ...form, password: e.target.value })}
                        className="w-full px-3 py-2 text-blackish border rounded focus:outline-none focus:ring   focus:bg-lightPurple focus:border-lightPurple bg-darkPurple"
                        required
                    />
                </div>

                <button className="w-full py-2 my-1 border border-lightPurple text-lightPurple hover:bg-midPurple hover:text-blackish active:bg-lightPurple rounded-sm transition">
                    Login
                </button>

            </form>
        </div>

    );
}
export default LoginForm;