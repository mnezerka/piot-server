import Home from './pages/Home.svelte'
import NotFound from './pages/NotFound.svelte'
import Customers from './pages/Customers.svelte'
import CustomerAdd from './pages/CustomerAdd.svelte'
import Users from './pages/Users.svelte'
import Devices from './pages/Devices.svelte'
import Login from './pages/Login.svelte'
import Signout from './pages/Signout.svelte'

let routes = {
    '/': Home,
    '/customers': Customers,
    '/customer-add': CustomerAdd,
    '/users': Users,
    '/devices': Devices,
    '/login': Login,
    '/signout': Signout,

    // Catch-all, must be last
    '*': NotFound,
}

export default routes
