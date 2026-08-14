import { Routes, Route } from 'react-router-dom'
import { Layout } from '@/components/layout/Layout'
import { DashboardPage } from '@/pages/DashboardPage'
import { HomePage } from '@/pages/HomePage'
import { SearchPage } from '@/pages/SearchPage'
import { BookPage } from '@/pages/BookPage'
import { CreateBookPage } from '@/pages/CreateBookPage'
import { LoginPage } from '@/pages/LoginPage'
import { RegisterPage } from '@/pages/RegisterPage'
import { ProfilePage } from '@/pages/ProfilePage'
import { UserProfilePage } from '@/pages/UserProfilePage'
import { ProtectedRoute } from '@/components/auth/ProtectedRoute'
import { RateLimitWarning } from '@/components/gateway/RateLimitWarning'
import { DevPanel } from '@/components/gateway/DevPanel'

function App() {
  return (
    <>
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<DashboardPage />} />
          <Route path="books" element={<HomePage />} />
          <Route path="search" element={<SearchPage />} />
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
      <RateLimitWarning />
      <DevPanel />
    </>
  )
}

export default App
