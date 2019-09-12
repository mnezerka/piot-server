import Home from './pages/Home.svelte'
import NotFound from './pages/NotFound.svelte'
import Users from './pages/Users.svelte'
import Login from './pages/Login.svelte'

let routes = {
    '/': Home,
    '/users': Users,
    '/login': Login,

    // Catch-all, must be last
    '*': NotFound,
}

export default routes
