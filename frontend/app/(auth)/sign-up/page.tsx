"use client";
import { CountrySelectField } from "@/components/forms/CountrySelectField";
import FooterLink from "@/components/forms/FooterLink";
import InputField from "@/components/forms/inputField";
import SelectField from "@/components/forms/SelectField";
import { Button } from "@/components/ui/button";
import { PREFERRED_INDUSTRIES } from "@/lib/constants";
import { useForm } from "react-hook-form";

const SignUp = () => {
	const {
		register,
		handleSubmit,
		control,
		formState: { errors, isSubmitting },
	} = useForm<SignUpFormData>({
		defaultValues: {
			fullName: "",
			email: "",
			password: "",
			country: "BG",
			preferredIndustry: "Technology",
		},
		mode: "onBlur",
	});
	const onSubmit = async (data: SignUpFormData) => {
		try {
			console.log(data);
		} catch (e) {
			console.error(e);
		}
	};

	return (
		<>
			<h1 className="form-title">Sign Up & Personalize</h1>
			<form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
				<InputField
					name="fullname"
					label="Full Name"
					placeholder="John Doe"
					register={register}
					error={errors.fullName}
					validation={{
						required: "Full name is required!",
						minLength: 2,
					}}
				/>
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
				<CountrySelectField
					name="country"
					label="Country"
					control={control}
					error={errors.country}
					required
				/>
				<SelectField
					name="preferredIndustry"
					label="Preferred Industry"
					placeholder="Select your preferred industry"
					options={PREFERRED_INDUSTRIES}
					control={control}
					error={errors.preferredIndustry}
					required
				/>
				<Button
					type="submit"
					disabled={isSubmitting}
					className="yellow-btn w-full mt-5"
				>
					{isSubmitting
						? "Creating Account"
						: "Start Your Investing Journey"}
				</Button>

				<FooterLink
					text="Already have an account?"
					linkText="Sign in"
					href="/sign-in"
				/>
			</form>
		</>
	);
};

export default SignUp;
