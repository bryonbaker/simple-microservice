using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using CounterClaim.Models;

namespace CounterClaim.Controllers
{
    public class APIController : Controller
    {
        public class Model
        {
            public int X { get; set; }
            public int Y { get; set; }
        }

        public class Result
        {
            public int Id { get; set; }
        }

        [HttpPost]
        [ResponseCache(Duration = 0, Location = ResponseCacheLocation.None, NoStore = true)]
        public JsonResult Count([FromBody]Model model)
        {
            return Json(new Result { Id = (model.X + 1) * (model.Y + 1) });
        }

        public class PersonModel
        {
            public string FirstName { get; set; }
            public string LastName { get; set; }
            public string Abn { get; set; }
        }

        public class PostModelResult
        {
            public bool ValidFirstName { get; set; }
            public bool ValidLastName { get; set; }
            public string AbnStatus { get; set; }
            public string Message { get; set; }
        }
        [HttpPost]
        [ResponseCache(Duration = 0, Location = ResponseCacheLocation.None, NoStore = true)]
        public JsonResult Form([FromBody]PersonModel model)
        {
            return Json(new PostModelResult { AbnStatus = "ABN is valid", ValidFirstName = true, ValidLastName = true, Message = "It works!" });
        }
    }
}
