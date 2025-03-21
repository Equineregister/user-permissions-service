# user_permissions_service - Tenant Specific Data Population


This directory contains a data population SQL script for each Tenant.
A script must exist for each Tenant. If one Tenant has the exact same data population requirements of another Tenant then copy the script - don't share it!

The data population scripts are intended as a foundation for the contents of the Tenant's database, regular API use will continue population of the database.

No consideration has been made for revert scripts.

For Tenant IDs, refer to https://scanimal.atlassian.net/wiki/spaces/WG/pages/1590820882/Tenant+IDs

