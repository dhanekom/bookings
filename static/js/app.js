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
      icon: icon,
      position: position,
      showConfirmButton: false,
      timer: 3000,
      timerProgressBar: true,
      didOpen: (toast) => {
        toast.addEventListener('mouseenter', Swal.stopTimer)
        toast.addEventListener('mouseleave', Swal.resumeTimer)
      }
    })

    Toast.fire({})
  }

  let success = function (c) {
    const {
      msg = "",
      title = "",
      footer = "",
    } = c;

    Swal.fire({
      icon: 'success',
      title: title,
      text: msg,
      footer: footer,
    })
  }

  let error = function (c) {
    const {
      msg = "",
      title = "",
      footer = "",
    } = c;

    Swal.fire({
      icon: 'error',
      title: title,
      text: msg,
      footer: footer,
    })
  }

  async function custom(c) {
    const {
      icon = "",
      msg = "",
      title = "",
      showConfirmButton = true,
    } = c;

    const { value: result } = await Swal.fire({
      icon: icon,
      title: title,
      html: msg,
      backdrop: false,
      focusConfirm: false,
      showCancelButton: true,
      showConfirmButton: showConfirmButton,
      willOpen: () => {
        if (c.willOpen !== undefined) {
          c.willOpen();
        }
      },
      preConfirm: () => {
        if (c.preConfirm !== undefined) {
          c.preConfirm();
        }
      },
      didOpen: () => {
        if (c.didOpen !== undefined) {
          c.didOpen();
        }
      }
    })

    if (result) {
      if (result.dismiss !== Swal.DismissReason.cancel && result.value !== "") {
        c.callback(result);
      } else {
        c.callback(false);
      }
    }
  }

  return {
    toast: toast,
    success: success,
    error: error,
    custom: custom,
  }
}

function addCheckAvailability(roomID, CSRFToken){
  document.getElementById("check-availability-button").addEventListener("click", function () {
    let html = `
        <form id="check-availability-form" action="" method="post" novalidate class="needs-validation">
          <div class="row mb-2" id="reservation-dates-modal">
            <div class="col">
              <input disabled required class="form-control" type="text" name="start" id="start" placeholder="Arrival date" autocomplete="off">
            </div>
            <div class="col">
              <input disabled required class="form-control" type="text" name="end" id="end" placeholder="Departure date" autocomplete="off">
            </div>
          </div>
        </form>
      `;
    attention.custom({
      title: "Choose your dates",
      msg: html,
      willOpen: () => {
        const elem = document.getElementById('reservation-dates-modal');
        const rp = new DateRangePicker(elem, {
          format: 'yyyy-mm-dd',
          showOnFocus: true,
          minDate: new Date(),
        });
      },
      didOpen: () => {
        document.getElementById('start').removeAttribute('disabled');
        document.getElementById('end').removeAttribute('disabled');
      },
      callback: function (result) {
        console.log(result);

        let form = document.getElementById("check-availability-form")
        let formData = new FormData(form)
        formData.append("csrf_token", CSRFToken)
        formData.append("room_id", roomID)


        fetch('/search-availability-json', {
          method: "post",
          body: formData,
        })
          .then(response => response.json())
          .then(data => {
            if (data.ok) {
              attention.custom({
                icon: 'success',
                showConfirmButton: false,
                msg: '<p>Room is available!</p>'
                    +'<p><a href="/book-room?id='+data.room_id+'&s='+data.start_date+'&e='+data.end_date+'"'
                    +' class="btn btn-primary">Book now!</a></p>'
              })
            } else {
              attention.error({
                msg: "Room is not available for the selected dates"
              })
            }
          })
      }
    })
  });    
}