$(function () {
    console.log('form started');
    $(".hide-me").hide();
    $("form").on("submit", function (e) {
        e.preventDefault();

        console.log("form submitting");

        let form = $(this);
        $.ajax({
            type: "post",
            url: "/api/form",
            contentType: "application/json",
            data: form.serialize(),
            success: function (data) {
                $("#_firstNameValid").text("First Name Valid: " + data.validFirstName);
                $("#_lastNameValid").text("Last Name Valid: " + data.validLastName);
                $("#_abnStatus").text("ABN Status: " + data.abnStatus);
                $("#_message").text("Message: " + data.message);
                $(".hide-me").show();
            }
        });
    });
});