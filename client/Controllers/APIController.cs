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
    }
}
