import { Routes, Route } from 'react-router-dom'
import { Layout } from '@/components/layout/Layout.tsx'
import { HomePage } from '@/pages/HomePage.tsx'
import { BookPage } from '@/pages/BookPage.tsx'
import { CreateBookPage } from '@/pages/CreateBookPage.tsx'
import { LoginPage } from '@/pages/LoginPage.tsx'
import { RegisterPage } from '@/pages/RegisterPage.tsx'
import { ProfilePage } from '@/pages/ProfilePage.tsx'
import { UserProfilePage } from '@/pages/UserProfilePage.tsx'
import { ProtectedRoute } from '@/components/auth/ProtectedRoute.tsx'

function App() {
  return (
    <Routes>
      <Route path="/" element={<Layout />}>
        <Route index element={<HomePage />} />
        <Route path="books/:id" element={<BookPage />} />
        <Route
          path="books/new"
          element={
            <ProtectedRoute>
              <CreateBookPage />
            </ProtectedRoute>
          }
        />
        <Route path="login" element={<LoginPage />} />
        <Route path="register" element={<RegisterPage />} />
        <Route path="users/:userId" element={<UserProfilePage />} />
        <Route
          path="profile"
          element={
            <ProtectedRoute>
              <ProfilePage />
            </ProtectedRoute>
          }
        />
      </Route>
    </Routes>
  )
}

export default App

