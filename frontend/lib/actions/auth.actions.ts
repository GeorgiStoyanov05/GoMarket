'use server';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000';

export async function getCurrentUser(): Promise<User | null> {
  try {
    const response = await fetch(`${API_URL}/auth/me`, {
      method: 'GET',
      credentials: 'include',
      cache: 'no-store',
    });

    if (!response.ok) {
      // User is not authenticated
      return null;
    }

    const data = await response.json();
    
    // Map the backend response to the frontend User type
    // The /auth/me endpoint now returns {user: {id, fullName, email, ...}}
    const userData = data.user || data;
    
    return {
      id: userData.id || '',
      name: userData.fullName || userData.name || userData.email?.split('@')[0] || 'User',
      email: userData.email || '',
    };
  } catch (error) {
    console.error('Error fetching current user:', error);
    return null;
  }
}
