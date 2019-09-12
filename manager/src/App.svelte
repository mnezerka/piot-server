<script>
    import Router from 'svelte-spa-router';
    import {link} from 'svelte-spa-router';
    import routes from './routes.js';
    import Navbar from './components/Navbar.svelte';
    import {token, authenticated} from './stores.js'

    console.log($authenticated, $token);

</script>

<Navbar/>

<main class="piot-main">
    <Router {routes}/>
</main>
