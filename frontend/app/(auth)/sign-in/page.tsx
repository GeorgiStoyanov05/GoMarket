"use client";
import FooterLink from "@/components/forms/FooterLink";
import InputField from "@/components/forms/InputField";
import { Button } from "@/components/ui/button";
import { useForm } from "react-hook-form";
import { useRouter } from "next/navigation";

const SignIn = () => {
	const {
		register,
		handleSubmit,
		formState: { errors, isSubmitting },
	} = useForm<SignInFormData>({
		defaultValues: {
			email: "",
			password: "",
		},
		mode: "onBlur",
	});

	const router = useRouter();

	const onSubmit = async (data: SignInFormData) => {
		const res = await fetch(
			`${process.env.NEXT_PUBLIC_API_URL}/auth/login`,
			{
				method: "POST",
				headers: { "Content-Type": "application/json" },
				credentials: "include",
				body: JSON.stringify(data),
			},
		);

		if (!res.ok) {
			const body = await res.json().catch(() => null);
			throw new Error(body?.error ?? "Login failed");
		}
		if (res.ok) router.replace("/");
	};
	return (
		<>
			<h1 className="form-title">Log In Your Account</h1>
			<form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
				<InputField
					name="email"
					label="Email"
					placeholder="Enter your email"
					register={register}
					error={errors.email}
					validation={{
						required: "Email is required!",
						pattern: /^\w+@\w+\.\w+$/,
						message: "Email is not in correct format!",
					}}
				/>
				<InputField
					name="password"
					label="Password"
					placeholder="Enter a strong password"
					type="password"
					register={register}
					error={errors.password}
					validation={{
						required: "Password is required!",
						minLength: 8,
					}}
				/>
				<Button
					type="submit"
					disabled={isSubmitting}
					className="yellow-btn w-full mt-5"
				>
					{isSubmitting ? "Logging in" : "Log in"}
				</Button>

				<FooterLink
					text="Don't have an account?"
					linkText="Sign up"
					href="/sign-up"
				/>
			</form>
		</>
	);
};

export default SignIn;
