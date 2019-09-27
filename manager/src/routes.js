import Home from './pages/Home.svelte'
import NotFound from './pages/NotFound.svelte'
import Orgs from './pages/Orgs.svelte'
import OrgView from './pages/OrgView.svelte'
import OrgAdd from './pages/OrgAdd.svelte'
import OrgEdit from './pages/OrgEdit.svelte'
import Users from './pages/Users.svelte'
import User from './pages/User.svelte'
import Things from './pages/Things.svelte'
import Login from './pages/Login.svelte'
import Signout from './pages/Signout.svelte'

let routes = {
    '/': Home,
    '/orgs': Orgs,
    '/org-view/:id': OrgView,
    '/org-add': OrgAdd,
    '/org-edit/:id': OrgEdit,
    '/users': Users,
    '/user/:id': User,
    '/things': Things,
    '/login': Login,
    '/signout': Signout,

    // Catch-all, must be last
    '*': NotFound,
}

export default routes
