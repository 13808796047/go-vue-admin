$ErrorActionPreference = "Stop"
$B = "http://127.0.0.1:8080"

function PostApi($url, $body, $tok) {
    $h = @{}
    if ($tok) { $h["Authorization"] = "Bearer $tok" }
    return (Invoke-RestMethod -Method Post -Uri $url -Headers $h -ContentType "application/json" -Body ($body | ConvertTo-Json -Depth 5))
}
function GetApi($url, $tok) {
    $h = @{}
    if ($tok) { $h["Authorization"] = "Bearer $tok" }
    return (Invoke-RestMethod -Method Get -Uri $url -Headers $h)
}
function PutApi($url, $body, $tok) {
    $h = @{}
    if ($tok) { $h["Authorization"] = "Bearer $tok" }
    return (Invoke-RestMethod -Method Put -Uri $url -Headers $h -ContentType "application/json" -Body ($body | ConvertTo-Json -Depth 5))
}
function DelApi($url, $tok) {
    $h = @{ "Authorization" = "Bearer $tok" }
    return (Invoke-RestMethod -Method Delete -Uri $url -Headers $h)
}
function Catch($script) {
    try { & $script } catch {
        $r = $_.Exception.Response
        if ($r) {
            $sr = [System.IO.StreamReader]::new($r.GetResponseStream())
            $txt = $sr.ReadToEnd()
            return ($txt | ConvertFrom-Json)
        }
        return @{ code = -1; msg = $_.Exception.Message }
    }
}

Write-Host "== 1. login bad password -> 401 =="
$r = Catch { PostApi "$B/api/login" @{ username = "admin"; password = "wrong" } $null }
Write-Host ("  code=" + $r.code + " msg=" + $r.msg)

Write-Host "== 2. login ok (admin/123456) =="
$r = PostApi "$B/api/login" @{ username = "admin"; password = "123456" } $null
$TOKEN = $r.data.accessToken
$REFRESH = $r.data.refreshToken
$roleList = $r.data.user.roles | ForEach-Object { $_.code }
Write-Host ("  user=" + $r.data.user.username + " roles=" + ($roleList -join ","))

Write-Host "== 3. no auth -> 401 =="
$r = Catch { GetApi "$B/api/users" $null }
Write-Host ("  code=" + $r.code + " msg=" + $r.msg)

Write-Host "== 4. bad token -> 401 =="
$r = Catch { GetApi "$B/api/users" "bogus" }
Write-Host ("  code=" + $r.code + " msg=" + $r.msg)

Write-Host "== 5. refresh =="
$r = PostApi "$B/api/refresh" @{ refreshToken = $REFRESH } $null
Write-Host ("  newAccessToken=" + $r.data.accessToken.Substring(0, 20) + "...")

Write-Host "== 6. /me =="
$r = GetApi "$B/api/me" $TOKEN
Write-Host ("  me=" + $r.data.username)

Write-Host "== 7. users page =="
$r = GetApi "$B/api/users?page=1&pageSize=10" $TOKEN
Write-Host ("  total=" + $r.data.total + " first=" + $r.data.list[0].username)

Write-Host "== 8. create role operator (need id before user) =="
$r = PostApi "$B/api/roles" @{ name = "operator"; code = "operator"; sort = 1; remark = "ops"; status = 1 } $TOKEN
$OPROLE = $r.data.id
Write-Host ("  role id=" + $OPROLE)

Write-Host "== 9. create user alice bound to operator =="
$r = PostApi "$B/api/users" @{ username = "alice"; nickname = "Alice"; password = "alice123"; email = "a@x.com"; roleIds = @($OPROLE); status = 1 } $TOKEN
Write-Host ("  alice id=" + $r.data.id)

Write-Host "== 10. user detail shows roles =="
$r = GetApi "$B/api/users/2" $TOKEN
Write-Host ("  nick=" + $r.data.nickname + " roles=" + (($r.data.roles | ForEach-Object { $_.code }) -join ","))

Write-Host "== 11. update user nickname =="
$r = PutApi "$B/api/users/2" @{ nickname = "Alice2"; email = "a2@x.com"; roleIds = @($OPROLE) } $TOKEN
Write-Host ("  nick=" + $r.data.nickname)

Write-Host "== 12. roles page =="
$r = GetApi "$B/api/roles" $TOKEN
Write-Host ("  total=" + $r.data.total + " first=" + $r.data.list[0].name)

Write-Host "== 13. assign menus+apis to operator =="
$r = PutApi "$B/api/roles/$OPROLE/menus" @{ menuIds = @(1,2,5) } $TOKEN
Write-Host ("  menus assigned code=" + $r.code)
$r = PutApi "$B/api/roles/$OPROLE/apis" @{ apiIds = @(6,7) } $TOKEN
Write-Host ("  apis assigned code=" + $r.code)

Write-Host "== 14. menu tree =="
$r = GetApi "$B/api/menus" $TOKEN
$top = ($r.data | ForEach-Object { $_.name }) -join ","
$child = (($r.data[0].children | ForEach-Object { $_.name })) -join ","
Write-Host ("  roots=" + $top + " under system=" + $child)

Write-Host "== 15. admin permissions =="
$r = GetApi "$B/api/permissions" $TOKEN
Write-Host ("  menus=" + $r.data.menus.Count + " apis=" + $r.data.apis.Count + " roleCodes=" + ($r.data.roleCodes -join ","))

Write-Host "== 16. alice permissions (operator only) =="
$r = PostApi "$B/api/login" @{ username = "alice"; password = "alice123" } $null
$ATOKEN = $r.data.accessToken
$r = GetApi "$B/api/permissions" $ATOKEN
Write-Host ("  alice menus=" + $r.data.menus.Count + " apis=" + $r.data.apis.Count + " roleCodes=" + ($r.data.roleCodes -join ","))

Write-Host "== 17. alice allowed: GET /api/users =="
$r = Catch { GetApi "$B/api/users" $ATOKEN }
Write-Host ("  code=" + $r.code + " (expect 200)")

Write-Host "== 18. alice forbidden: POST /api/roles =="
$r = Catch { PostApi "$B/api/roles" @{ name = "X"; code = "xx" } $ATOKEN }
Write-Host ("  code=" + $r.code + " msg=" + $r.msg)

Write-Host "== 19. change password =="
$r = PutApi "$B/api/me/password" @{ oldPassword = "alice123"; newPassword = "alice456" } $ATOKEN
Write-Host ("  code=" + $r.code)

Write-Host "== 20. old password rejected -> 401 =="
$r = Catch { PutApi "$B/api/me/password" @{ oldPassword = "alice123"; newPassword = "x12345" } $ATOKEN }
Write-Host ("  code=" + $r.code + " msg=" + $r.msg)

Write-Host "== 21. update profile =="
$r = PutApi "$B/api/profile" @{ nickname = "Super"; email = "root@x.com" } $TOKEN
Write-Host ("  nick=" + $r.data.nickname + " email=" + $r.data.email)

Write-Host "== 22. soft delete user 2 =="
$r = DelApi "$B/api/users/2" $TOKEN
Write-Host ("  code=" + $r.code)
$r = Catch { GetApi "$B/api/users/2" $TOKEN }
Write-Host ("  after: code=" + $r.code + " msg=" + $r.msg)

Write-Host "== 23. menu with child; delete parent must fail =="
$r = PostApi "$B/api/menus" @{ parentId = 0; name = "report"; permKey = "report"; type = 2; icon = "chart"; url = "/report"; component = "report/Index"; sort = 3; status = 1 } $TOKEN
$MID = $r.data.id
Write-Host ("  created menu id=" + $MID)
$r = PostApi "$B/api/menus" @{ parentId = $MID; name = "sale"; permKey = "report:sale"; type = 3 } $TOKEN
Write-Host ("  child created code=" + $r.code)
$r = Catch { DelApi "$B/api/menus/$MID" $TOKEN }
Write-Host ("  delete parent: code=" + $r.code + " msg=" + $r.msg)

Write-Host "== 24. duplicate username rejected =="
$r = Catch { PostApi "$B/api/users" @{ username = "admin"; nickname = "x"; password = "x123456" } $TOKEN }
Write-Host ("  code=" + $r.code + " msg=" + $r.msg)

Write-Host "== 25. param validation (short username+password) -> 400 =="
$r = Catch { PostApi "$B/api/users" @{ username = "b"; nickname = "x"; password = "123" } $TOKEN }
Write-Host ("  code=" + $r.code + " msg=" + $r.msg)

Write-Host "== 26. create bob, disable, login blocked =="
$r = PostApi "$B/api/users" @{ username = "bob"; nickname = "Bob"; password = "bob12345"; roleIds = @(); status = 1 } $TOKEN
$BOBID = $r.data.id
Write-Host ("  bob id=" + $BOBID)
$r = PutApi "$B/api/users/$BOBID/status" @{ status = 0 } $TOKEN
Write-Host ("  disable code=" + $r.code)
$r = Catch { PostApi "$B/api/login" @{ username = "bob"; password = "bob12345" } $null }
Write-Host ("  login: code=" + $r.code + " msg=" + $r.msg)

Write-Host "== 27. delete operator role =="
$r = DelApi "$B/api/roles/$OPROLE" $TOKEN
Write-Host ("  code=" + $r.code)

Write-Host "== 28. delete super admin role blocked =="
$r = Catch { DelApi "$B/api/roles/1" $TOKEN }
Write-Host ("  code=" + $r.code + " msg=" + $r.msg)

Write-Host "SMOKE DONE"
