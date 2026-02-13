import { BrowserRouter, Routes, Route, Outlet, NavLink } from 'react-router'

import Search from './Search.tsx'

export function Layout() {
  return (
    <div className="flex flex-col">
      <div className="flex flex-row justify-between flex-1 bg-slate-500 p-4">
        <NavLink className="text-lg" to="/">
          pictag
        </NavLink>
        <NavLink className="btn" to="/add">
          Add Image
        </NavLink>
      </div>
      <div>
        <Outlet />
      </div>
    </div>
  )
}

export function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<Layout />}>
          <Route index element={<Search />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}

export default App
