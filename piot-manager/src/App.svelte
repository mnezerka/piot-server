<script>
    import Router from 'svelte-spa-router';
    import {link} from 'svelte-spa-router';
    import routes from './routes.js';
    import Navbar from './components/Navbar.svelte';
    import {profile, token, authenticated} from './stores.js'
    import {onMount} from 'svelte';
    import {gql} from './utils';

    var authenticating = true;

    onMount(async () => {

        // verify that user token is valid - download user profile
        if ($token) {

            try {
                let data = await gql({query: "{userProfile {email}}"});
                $authenticated = true;
                $profile = data.userProfile;
            } catch(error) {
                error = 'Request failed (' + error + ')';
            }

        }
        authenticating = false;
    })

</script>

<Navbar/>

<main class="piot-main">
    <div class="container">
        {#if authenticating}
        <span>Authenticating</span>
        {:else}
        <Router {routes}/>
        {/if}
    </div>
</main>
