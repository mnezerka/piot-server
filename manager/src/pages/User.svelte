<script>
    import {token, authenticated} from '../stores.js'
    import {push} from 'svelte-spa-router'
    import {gql} from '../utils.js';
    import {onMount} from 'svelte';
    import ErrorBar from '../components/ErrorBar.svelte';

    export var params;

    let email = '';
    let user = null;
    let error = null;
    let fetching = false;
    let success = null;
    let orgs = [];
    let orgAdd = '';

    async function fetchUser() {
        if (!params.id) {
            error = 'No user specified';
            return
        }

        fetching = true;
        error = null

        try {
            let data = await gql({query: `{user(id: "${params.id}") {id, email, created} orgs {id, name}}`});
            user = data.user;
            email = user.email;
            orgs = data.orgs;
        } catch (e) {
            error = e;
        }
        fetching = false;
    }

    async function updateUser()
    {
        if (name.length === 0) {
            error = 'No name specified'
            return
        }

        fetching = true;
        error = null;
        success = false;

        try {
            let data = await gql({query: `mutation {updateUser(user: {id: "${params.id}", email: "${email}"}) {id}}`});
            success = 'User successfully updated'
        } catch(e) {
            error = e;
        }
        fetching = false;
    }

    onMount(fetchUser)

    async function addOrg() {
        console.log('adding org', orgAdd)
    }

</script>

<style>
form { width: 24rem;}
h2 { margin-top: 2rem; }
</style>

<h1 class="title">User</h1>

<ErrorBar error={error}/>

{#if success}<div class="notification is-success has-text-centered">{success}</div>{/if}

{#if user}
<form on:submit|preventDefault={updateUser}>

    <div class="field">
        <p class="control">
            <input bind:value={email} class="input {email.length === 0 && "is-danger"}" placeholder="User email">
        </p>
    </div>

    <div class="field">
        <p class="control">
            <button class="button is-success">Update</button>
        </p>
    </div>

</form>

<h2 class="subtitle">Organizations</h2>

<form on:submit|preventDefault={addOrg}>

    <div class="field is-grouped">
        <p class="control is-expanded">
            <input bind:value={orgAdd} class="input">
        </p>

         <p class="control">
             <button class="button is-success">Add</button>
         </p>
    </div>

</form>

{/if}
