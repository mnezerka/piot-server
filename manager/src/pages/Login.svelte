<script>
    import {token, authenticated} from '../stores.js'
    import {push} from 'svelte-spa-router'

    let username = '';
    let password = '';
    let error = null;

    async function handleSubmit()
    {
        let response;

        // fetch data from server
        try {
            response = await fetch('http://localhost:9096/login', {
                method: 'POST',
                mode: 'cors',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({email: username, password})
            })
        } catch(error) {
            error = 'Login Request Failed (' + error + ')';
        }

        // proces data
        try {
            let data = await response.json();

            // if the response status code != 200 OK
            if (!response.ok) {
                error = 'Login Request Failed'

                // try to add more information to the error message
                if (data.error) {
                    error = `${error} (${data.error})`;
                }
            } else if (data.token) {
                localStorage.setItem('token', data.token);
                $token = data.token;
                $authenticated = true;
                push('/');
            }
        } catch(error) {
            error = 'Login Request Failed (' + error + ')';
        }
    }
</script>

<style>
form {
    margin-top: 100px;
}
</style>

{#if error}
    <div class="notification is-danger has-text-centered">
        {error}
    </div>
{/if}

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
