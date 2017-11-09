namespace Swirl.Setting {
    export class IndexPage {
        constructor() {
            $("#ldap-enabled").change(e => {
                let enabled = $(e.target).prop("checked");
                $("#fs-ldap").find("input:not(:checkbox)").prop("readonly", !enabled);
            });
            $("#ldap-auth-simple,#ldap-auth-bind").click(e => {
                if ($(e.target).val() == "0") {
                    $("#div-auth-simple").show();
                    $("#div-auth-bind").hide();
                } else {
                    $("#div-auth-simple").hide();
                    $("#div-auth-bind").show();
                }
            });
        }
    }
}