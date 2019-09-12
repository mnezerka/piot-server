<script>
    import Router from 'svelte-spa-router';
    import {link} from 'svelte-spa-router';
    import routes from './routes.js';
    import Navbar from './components/Navbar.svelte';
    import {profile, token, authenticated} from './stores.js'

    // verify that user token is valid - download user profile
    if ($token) {
        var requestBody = {
            query: "{userProfile {email}}"}
        }
        fetch('http://localhost:9096/query', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + $token,
            },
            body: JSON.stringify(requestBody)
        })
            .then(function(response) {
                if (!response.ok) {
                    throw response
                }
                return response.json()
            })
            .then(function(data) {
                $authenticated = true;
                $profile = data.data.userProfile;
            })
            .catch(function(response) {
                console.log('Failed to fetch page: ', response);
            });

</script>

<Navbar/>

<main class="piot-main">
    <Router {routes}/>
</main>
