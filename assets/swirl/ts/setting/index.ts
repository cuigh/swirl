namespace Swirl.Setting {
    export class IndexPage {
        constructor() {
            $("#ldap-enabled").change(e => {
                let enabled = $(e.target).prop("checked");
                $("#fs-ldap").find("input:not(:checkbox)").prop("readonly", !enabled);
            });
        }
    }
}