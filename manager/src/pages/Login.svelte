<script>
    import {token, authenticated} from '../stores.js'
    import {push} from 'svelte-spa-router'

    let username = '';
    let password = '';

    function handleSubmit()
    {
        console.log('handle submit', username, password);

        fetch('http://localhost:9096/login', {
            method: 'POST',
            mode: 'cors',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({email: username, password})
        })
            .then(function(response) {
                if (!response.ok) {
                    throw response
                }
                return response.json()
            })
            .then(function(data) {
                console.log(data);
                if (data.token) {
                    localStorage.setItem('token', data.token);
                    $token = data.token;
                    $authenticated = true;
                    push('/');
                }
            })
            .catch(function(response) {
                console.log('Failed to fetch page: ', response);
            });


    }
</script>

<style>
form {
    margin-top: 100px;
}
</style>

<div class="columns is-mobile is-centered">
  <div class="column is-one-quarter is-vcentered">

<form class="has-text-centered" on:submit|preventDefault={handleSubmit}>

    <div class="field">
        <p class="control has-icons-left has-icons-right">
            <input bind:value={username} class="input" type="email" placeholder="Email">
            <span class="icon is-small is-left">
                <i class="fas fa-envelope"></i>
            </span>
            <span class="icon is-small is-right">
                <i class="fas fa-check"></i>
            </span>
        </p>
    </div>

    <div class="field">
        <p class="control has-icons-left">
            <input bind:value={password} class="input" type="password" placeholder="Password">
            <span class="icon is-small is-left">
                <i class="fas fa-lock"></i>
            </span>
        </p>
    </div>

    <div class="field">
        <p class="control">
            <button class="button is-success">Login</button>
        </p>
    </div>

</form>

  </div>
</div>
