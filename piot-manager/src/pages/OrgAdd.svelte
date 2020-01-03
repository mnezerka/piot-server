<script>
    import {token, authenticated} from '../stores.js'
    import {push} from 'svelte-spa-router'
    import {gql} from '../utils.js';
    import {onMount} from 'svelte';

    let name = '';
    let description = '';
    let error = null;
    let fetching = false;
    let success = null;

    onMount(() => {
        if (!$authenticated) { push("/login"); }
    })

    async function handleSubmit()
    {
        if (name.length === 0) {
            error = 'No name specified'
            return
        }

        fetching = true;
        error = false;
        success = false;

        try {
            let data = await gql({query: `mutation {createOrg(name: "${name}", description: "${description}") {name}}`});
            success = 'Organization successfully created'
        } catch(e) {
            if (e instanceof Array) {
                e = e.map((err) => err.message).join(', ');
            }
            error = 'Request failed (' + e + ')';
        }

        fetching = false;
    }
</script>

<style>
form { width: 24rem;}
</style>

<h1 class="title">Add Organization</h1>

{#if error}<div class="notification is-danger has-text-centered">{error}</div>{/if}
{#if success}<div class="notification is-success has-text-centered">{success}</div>{/if}

<form on:submit|preventDefault={handleSubmit}>

    <div class="field">
        <p class="control">
            <input bind:value={name} class="input {name.length === 0 && "is-danger"}" placeholder="Organization name">
        </p>
    </div>

    <div class="field">
        <p class="control">
            <textarea bind:value={description} class="textarea" placeholder="Organization description"/>
        </p>
    </div>

    <div class="field">
        <p class="control">
            <button class="button is-success">Create</button>
        </p>
    </div>

</form>

