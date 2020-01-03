<script>
    import {token, authenticated} from '../stores.js'
    import {push} from 'svelte-spa-router'
    import {gql} from '../utils.js';
    import {onMount} from 'svelte';
    import ErrorBar from '../components/ErrorBar.svelte';
    import TableButtonDelete from '../components/TableButtonDelete.svelte';

    export var params;

    let email = '';
    let user = null;
    let error = null;
    let fetching = false;
    let success = null;
    let orgs = [];
    let orgsAssigned = [];
    let orgAdd = '';

    onMount(() => {
        if (!$authenticated) { push("/login"); }
        fetchUser();
    })

    async function fetchUser() {
        if (!params.id) {
            error = 'No user specified';
            return
        }

        fetching = true;
        error = null

        try {
            let data = await gql({query: `{user(id: "${params.id}") {id, email, created, orgs {id, name}} orgs {id, name}}`});
            user = data.user;
            email = user.email;
            orgs = data.orgs;
            orgsAssigned = user.orgs.map(o => o.id);
        } catch (e) {
            error = e;
        }
        fetching = false;
    }

    async function updateUser()
    {
        if (email.length === 0) {
            error = 'No email specified'
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

    async function addOrg() {

        if (orgAdd.length === 0) {
            error = 'No organization specified'
            return
        }

        fetching = true;
        error = null;
        success = false;

        try {
            let data = await gql({query: `mutation {addOrgUser(orgId: "${orgAdd}", userId: "${params.id}")}`});
            success = 'User successfully added to organization'
        } catch(e) {
            error = e;
        }
        fetching = false;

        orgAdd = '';

        await fetchUser()
    }

    async function removeOrg(e) {

        if (!e.detail.value || e.detail.value.length === 0) {
            error = 'No organization specified'
            return
        }

        fetching = true;
        error = null;
        success = false;

        try {
            let data = await gql({query: `mutation {removeOrgUser(orgId: "${e.detail.value}", userId: "${params.id}")}`});
            success = 'User successfully removed from organization'
        } catch(e) {
            error = e;
        }
        fetching = false;

        await fetchUser()
    }

</script>

<style>
form { width: 24rem;}
h2 { margin-top: 2rem; }
.delete-button { text-align: right; }
</style>

<h1 class="title">User</h1>

<ErrorBar error={error}/>

{#if success}<div class="notification is-success has-text-centered">{success}</div>{/if}

{#if user}

<h2 class="subtitle">Edit User</h2>

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

    <div class="field has-addons">
        <div class="control is-expanded">
            <div class="select is-fullwidth">
                <select bind:value={orgAdd}>
                    <option value="">---</option>
                {#each orgs as org}
                    {#if orgsAssigned.indexOf(org.id) === -1}
                    <option value="{org.id}">{org.name}</option>
                    {/if}
                {/each}
                </select>
            </div>
        </div>
        <div class="control">
            <button class="button is-success">Add</button>
        </div>
    </div>

</form>

<table class="table is-fullwidth">
    <thead>
        <tr>
            <th>Name</th>
            <th></th>
        </tr>
    </thead>
    <tbody>
        {#each orgs as org}
            {#if orgsAssigned.indexOf(org.id) !== -1}
        <tr>
            <td>{org.name}</td>
            <td class="delete-button">
                <TableButtonDelete value={org.id} on:click={removeOrg} />
            </td>
        </tr>
            {/if}
        {/each}
    </tbody>
</table>

{/if}
