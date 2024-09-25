// Prompt is our JavaScript module for all alerts,notifications , and custom popup dialogs
function Prompt() {
    let toast = function (c) {
        const {
            msg = "",
            icon = "success",
            position = "top-end"
        } = c;
        const Toast = Swal.mixin({
            toast: true,
            title: msg,
            position: position,
            icon: icon,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.onmouseenter = Swal.stopTimer;
                toast.onmouseleave = Swal.resumeTimer;
            }
        });
        Toast.fire({});
    }

    let success = function (c) {
        const {
            msg = "",
            title = "",
            footer = ""
        } = c;
        Swal.fire({
            icon: "success",
            title: title,
            text: msg,
            footer: footer
        });
    }

    let error = function (c) {
        const {
            msg = "",
            title = "",
            footer = ""
        } = c;
        Swal.fire({
            icon: "error",
            title: title,
            text: msg,
            footer: footer
        });
    }

    async function custom(c) {
        const {
            icon = "",
            msg = "",
            title = "",
            showConfirmButton = true
        } = c;

        const {value: result} = await Swal.fire({
            icon:icon,
            title: title,
            html:msg,
            focusConfirm: false,
            backdrop:true,
            showCancelButton:true,
            showConfirmButton:showConfirmButton,
            willOpen:() => {
                if (c.willOpen !== undefined) {
                    c.willOpen()
                }
            },
            preConfirm: () => {
                return [
                    document.getElementById("start").value,
                    document.getElementById("end").value
                ];
            },
            didOpen:() => {
                if( c.didOpen !== undefined) {
                    c.didOpen();
                }
            }
        });
        if ( result ){
            // it's not because they clicked the cancel button.
            if (result.dismiss !== Swal.DismissReason.cancel) {
                if (result.value !== ""){
                    if(c.callback !== undefined) {
                        c.callback(result);
                    }
                }else {
                    c.callback(false);
                }
            }else {
                c.callback(false);
            }
        }
    }

    return {
        toast,
        success,
        error,
        custom:custom
    }
}


function CalenderPopup(roomID,CSRFToken) {
    document.getElementById("check-availability-button").addEventListener("click", function () {
        let html = `
        <form id="check-availability-form"  novalidate class="needs-validation">
            <div class="row">
                <div class="col">
                    <div class="row" id="reservation-dates-modal">
                        <div class="col">
                            <input disabled autocomplete="off" type="text" class="form-control" required name="start_date" id="start" placeholder="Arrival">
                        </div>
                        <div class="col">
                            <input disabled autocomplete="off" type="text" class="form-control" required name="end_date" id="end" placeholder="Departure">
                        </div>
                    </div>
                </div>
            </div>
        </form>
        `;
        attention.custom({
            msg:html,
            title:"Choose your dates.",
            willOpen: () => {
                const elem = document.getElementById("reservation-dates-modal");
                const rp = new DateRangePicker(elem,{
                    format:'yyyy-mm-dd',
                    showOnFocus:true,
                    minDate: new Date(),
                })

            },
            didOpen:() => {
                document.getElementById("start").removeAttribute('disabled');
                document.getElementById("end").removeAttribute('disabled');
            },
            callback:function (result) {

                let form = document.getElementById("check-availability-form");
                let formData = new FormData(form);
                formData.append("csrf_token",CSRFToken);
                formData.append("room_id",roomID);

                fetch('/search-availability-json',{
                    method:"post",
                    body: formData,

                })
                    .then(response => response.json())
                    .then(data => {
                        if (data.ok) {
                            attention.custom({
                                icon:"success",
                                msg: `<p>Room is available!</p>
                                    <p><a href='/book-room?id=${data.room_id}&s=${data.start_date}&e=${data.end_date}' class='btn btn-primary'>Book Now</a></p>`,
                                showConfirmButton:false
                            })
                        }else{
                            attention.error({
                                msg:"Room is Not availability!",
                            })
                        }
                    })
            }
        });
    })
}