<script>
    import {onMount} from 'svelte';
    import {authenticated, token} from '../stores'
    import OrgList from '../components/OrgList.svelte';
    import ErrorBar from '../components/ErrorBar.svelte';
    import {gql} from '../utils.js';
    import {push} from 'svelte-spa-router'

    let error = null;
    let orgs = null;
    let fetching = false;

    onMount(() => {
        if (!$authenticated) { push("/login"); }
        fetchOrgs();
    })

    async function fetchOrgs() {
        fetching = true;
        error = false;
        orgs = null;

        try {
            let data = await gql({query: "{orgs {id, name, created}}"});
            orgs = data.orgs;
        } catch(e) {
            error = e;//l'Request failed (' + e + ')';
        }

        fetching = false;
    }

    function onAdd() {
        push('/org-add');
    }

</script>


<h1 class="title">Organizations</h1>

{#if fetching}
    <progress class="progress is-small is-primary" max="100">15%</progress>
{:else}
    <ErrorBar error={error}/>
    {#if !error}
        <OrgList orgs={orgs}/>
        <button class="button" on:click={onAdd}>Add</button>
    {/if}
{/if}
