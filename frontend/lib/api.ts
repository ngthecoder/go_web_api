export async function updateProfile(token: string, username: string, email: string) {
  const response = await fetch('http://localhost:8000/api/user/profile/update', {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({ username, email })
  });

  if (!response.ok) {
    const error = await response.text();
    throw new Error(error);
  }

  return response.json();
}

export async function changePassword(token: string, currentPassword: string, newPassword: string) {
  const response = await fetch('http://localhost:8000/api/user/password', {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({ 
      current_password: currentPassword, 
      new_password: newPassword 
    })
  });

  if (!response.ok) {
    const error = await response.text();
    throw new Error(error);
  }

  return response.json();
}

export async function deleteAccount(token: string, password: string) {
  const response = await fetch('http://localhost:8000/api/user/account', {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({ password })
  });

  if (!response.ok) {
    const error = await response.text();
    throw new Error(error);
  }

  return response.json();
}
