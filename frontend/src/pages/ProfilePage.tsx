import { useState, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { User, Mail, Calendar, Shield, Save, Edit } from 'lucide-react';
import { useAuthStore } from '@/store/authStore';

interface ProfileFormData {
	username: string;
	email: string;
}

const ProfilePage = () => {
	const { user, updateUser } = useAuthStore();
	const [isEditing, setIsEditing] = useState(false);
	const [isSaving, setIsSaving] = useState(false);
	const [message, setMessage] = useState('');

	const {
		register,
		handleSubmit,
		formState: { errors },
		reset,
	} = useForm<ProfileFormData>();

	useEffect(() => {
		if (user) {
			reset({
				username: user.username,
				email: user.email,
			});
		}
	}, [user, reset]);

	const onSubmit = async (data: ProfileFormData) => {
		setIsSaving(true);
		setMessage('');

		try {
			// Simulate API call
			await new Promise(resolve => setTimeout(resolve, 1000));
			
			if (user) {
				updateUser({
					...user,
					username: data.username,
					email: data.email,
				});
			}
			
			setMessage('Profile updated successfully!');
			setIsEditing(false);
		} catch (error) {
			setMessage('Failed to update profile. Please try again.');
		} finally {
			setIsSaving(false);
		}
	};

	const handleCancel = () => {
		reset();
		setIsEditing(false);
		setMessage('');
	};

	if (!user) {
		return (
			<div className="text-center py-12">
				<p className="text-gray-600">Please log in to view your profile.</p>
			</div>
		);
	}

	return (
		<div className="max-w-4xl mx-auto space-y-6">
			{/* Header */}
			<div className="flex justify-between items-center">
				<div>
					<h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100">Profile</h1>
					<p className="text-gray-600 mt-1 dark:text-gray-300">
						Manage your account information and preferences
					</p>
				</div>
				<button
					onClick={() => setIsEditing(!isEditing)}
					className="bg-gray-200 text-gray-800 px-4 py-2 rounded-lg font-medium hover:bg-gray-300 transition-colors inline-flex items-center dark:bg-gray-800 dark:text-gray-200 dark:hover:bg-gray-700"
				>
					{isEditing ? (
						<>
							<Edit className="mr-2 h-4 w-4" />
							Cancel Edit
						</>
					) : (
						<>
							<Edit className="mr-2 h-4 w-4" />
							Edit Profile
						</>
					)}
				</button>
			</div>

			{/* Success Message */}
			{message && (
				<div className="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded-md">
					{message}
				</div>
			)}

			<div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
				{/* Profile Information */}
				<div className="lg:col-span-2">
					<div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6 dark:bg-gray-900 dark:border-gray-800">
						<h2 className="text-xl font-semibold text-gray-900 mb-6 dark:text-gray-100">Personal Information</h2>
						
						<form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
							<div className="grid grid-cols-1 md:grid-cols-2 gap-6">
								<div>
									<label htmlFor="username" className="block text-sm font-medium text-gray-700 mb-2 dark:text-gray-300">
										Username
									</label>
									<div className="relative">
										<input
											id="username"
											type="text"
											disabled={!isEditing}
											className={`w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent ${!isEditing ? 'bg-gray-50 cursor-not-allowed dark:bg-gray-800' : 'dark:bg-gray-900'} dark:border-gray-700 dark:text-gray-100`}
											{...register('username', {
												required: 'Username is required',
												minLength: {
													value: 3,
													message: 'Username must be at least 3 characters',
												},
											})}
										/>
										<User className="absolute right-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400" />
									</div>
									{errors.username && (
										<p className="mt-1 text-sm text-red-600">{errors.username.message}</p>
									)}
								</div>

								<div>
									<label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-2 dark:text-gray-300">
										Email Address
									</label>
									<div className="relative">
										<input
											id="email"
											type="email"
											disabled={!isEditing}
											className={`w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent ${!isEditing ? 'bg-gray-50 cursor-not-allowed dark:bg-gray-800' : 'dark:bg-gray-900'} dark:border-gray-700 dark:text-gray-100`}
											{...register('email', {
												required: 'Email is required',
												pattern: {
													value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
													message: 'Invalid email address',
												},
											})}
										/>
										<Mail className="absolute right-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400" />
									</div>
									{errors.email && (
										<p className="mt-1 text-sm text-red-600">{errors.email.message}</p>
									)}
								</div>
							</div>

							{/* Action Buttons */}
							{isEditing && (
								<div className="flex space-x-3 pt-4">
									<button
										type="submit"
										disabled={isSaving}
										className="bg-blue-600 text-white px-4 py-2 rounded-lg font-medium hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed inline-flex items-center dark:bg-blue-500 dark:hover:bg-blue-600"
									>
										{isSaving ? (
											<>
												<div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
												Saving...
											</>
										) : (
											<>
												<Save className="mr-2 h-4 w-4" />
												Save Changes
											</>
										)}
									</button>
									<button
										type="button"
										onClick={handleCancel}
										className="bg-gray-200 text-gray-800 px-4 py-2 rounded-lg font-medium hover:bg-gray-300 transition-colors dark:bg-gray-800 dark:text-gray-200 dark:hover:bg-gray-700"
									>
										Cancel
									</button>
								</div>
							)}
						</form>
					</div>
				</div>

				{/* Account Details */}
				<div className="lg:col-span-1">
					<div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6 dark:bg-gray-900 dark:border-gray-800">
						<h2 className="text-xl font-semibold text-gray-900 mb-6 dark:text-gray-100">Account Details</h2>
						
						<div className="space-y-4">
							<div className="flex items-center space-x-3">
								<div className="bg-blue-100 rounded-full p-2 dark:bg-blue-900/40">
									<Calendar className="h-5 w-5 text-blue-600" />
								</div>
								<div>
									<p className="text-sm text-gray-500 dark:text-gray-400">Member Since</p>
									<p className="font-medium text-gray-900 dark:text-gray-100">
										{new Date(user.created_at).toLocaleDateString()}
									</p>
								</div>
							</div>

							<div className="flex items-center space-x-3">
								<div className="bg-blue-100 rounded-full p-2 dark:bg-blue-900/40">
									<Shield className="h-5 w-5 text-blue-600" />
								</div>
								<div>
									<p className="text-sm text-gray-500 dark:text-gray-400">Roles</p>
									<div className="flex flex-wrap gap-1 mt-1">
										{user.roles.map((role) => (
											<span
												key={role.id}
												className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900/40 dark:text-blue-300"
											>
												{role.name}
											</span>
										))}
									</div>
								</div>
							</div>
						</div>

						{/* Quick Actions */}
						<div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6 mt-6 dark:bg-gray-900 dark:border-gray-800">
							<h3 className="text-lg font-semibold text-gray-900 mb-4 dark:text-gray-100">Quick Actions</h3>
							<div className="space-y-3">
								<button className="w-full text-left px-3 py-2 text-sm text-gray-700 hover:bg-gray-100 rounded-md transition-colors dark:text-gray-300 dark:hover:bg-gray-800">
									Change Password
								</button>
								<button className="w-full text-left px-3 py-2 text-sm text-gray-700 hover:bg-gray-100 rounded-md transition-colors dark:text-gray-300 dark:hover:bg-gray-800">
									Notification Settings
								</button>
								<button className="w-full text-left px-3 py-2 text-sm text-gray-700 hover:bg-gray-100 rounded-md transition-colors dark:text-gray-300 dark:hover:bg-gray-800">
									Privacy Settings
								</button>
							</div>
						</div>
					</div>
				</div>
			</div>
			</div>
		);
};

export default ProfilePage;

