-module(cfg_test_svr).

-export[get/1, get_id/1, get_name/1, get_age/1, get_sex/1, get_items/1].

-include("../include/cfg.hrl").

get(1) ->
	#cfg_test_svr{id=1,name="name1",age=10,sex=1,items=[1,2,3]};
get(2) ->
	#cfg_test_svr{id=2,name="name2",age=11,sex=2,items=[1,2,4]};
get(_) ->
	undefined.

get_id(Val=#cfg_test_svr{}) ->
	Val#cfg_test_svr.id;
get_id(_) ->
	undefined.

get_name(Val=#cfg_test_svr{}) ->
	Val#cfg_test_svr.name;
get_name(_) ->
	undefined.

get_age(Val=#cfg_test_svr{}) ->
	Val#cfg_test_svr.age;
get_age(_) ->
	undefined.

get_sex(Val=#cfg_test_svr{}) ->
	Val#cfg_test_svr.sex;
get_sex(_) ->
	undefined.

get_items(Val=#cfg_test_svr{}) ->
	Val#cfg_test_svr.items;
get_items(_) ->
	undefined.


