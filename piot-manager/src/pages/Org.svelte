<script>
    import {token, authenticated} from '../stores.js'
    import {push} from 'svelte-spa-router'
    import {gql} from '../utils.js';
    import {onMount} from 'svelte';
    import ErrorBar from '../components/ErrorBar.svelte';
    import TableButtonDelete from '../components/TableButtonDelete.svelte';

    export var params;

    let name = '';
    let org = null;
    let error = null;
    let fetching = false;
    let success = null;

    onMount(() => {
        if (!$authenticated) { push("/login"); }
        fetchOrg()
    })

    async function fetchOrg() {
        if (!params.id) {
            error = 'No organization specified';
            return
        }

        fetching = true;
        error = null

        try {
            let data = await gql({query: `{org(id: "${params.id}") {id, name, created, users {id, email}}}`});
            org = data.org;
            name = org.name;
        } catch (e) {
            error = e;
        }
        fetching = false;
    }

    async function updateOrg()
    {
        if (name.length === 0) {
            error = 'No name specified'
            return
        }

        fetching = true;
        error = null;
        success = false;

        try {
            let data = await gql({query: `mutation {updateOrg(org: {id: "${params.id}", name: "${name}"}) {id}}`});
            success = 'Organization successfully updated'
        } catch(e) {
            error = e;
        }
        fetching = false;
    }

    async function removeUser(e) {

        if (!e.detail.value || e.detail.value.length === 0) {
            error = 'No user specified'
            return
        }

        fetching = true;
        error = null;
        success = false;

        try {
            let data = await gql({query: `mutation {removeOrgUser(orgId: "${params.id}", userId: "${e.detail.value}")}`});
            success = 'User successfully removed from organization'
        } catch(e) {
            error = e;
        }
        fetching = false;

        await fetchOrg()
    }

</script>

<style>
form { width: 24rem;}
h2 { margin-top: 2rem; }
.delete-button { text-align: right; }
</style>

<h1 class="title">Organization</h1>

<ErrorBar error={error}/>

{#if success}<div class="notification is-success has-text-centered">{success}</div>{/if}

{#if org}

<h2 class="subtitle">Edit Organization</h2>

<form on:submit|preventDefault={updateOrg}>

    <div class="field">
        <p class="control">
            <input bind:value={name} class="input {name.length === 0 && "is-danger"}" placeholder="Organization name">
        </p>
    </div>

    <div class="field">
        <p class="control">
            <button class="button is-success">Update</button>
        </p>
    </div>

</form>

<h2 class="subtitle">Users</h2>

<table class="table is-fullwidth">
    <thead>
        <tr>
            <th>Email</th>
            <th></th>
        </tr>
    </thead>
    <tbody>
        {#each org.users as user}
        <tr>
            <td>{user.email}</td>
            <td class="delete-button">
                <TableButtonDelete value={user.id} on:click={removeUser} />
            </td>
        </tr>
        {/each}
    </tbody>
</table>

{/if}
