//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*

Copyright DNITSCH - WTFPL :)

Generated types for seeder used in apis module.

*/

// Code generated by deepcopy-gen. DO NOT EDIT.

package seeder

import (
	http "net/http"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Action) DeepCopyInto(out *Action) {
	*out = *in
	if in.header != nil {
		in, out := &in.header, &out.header
		*out = new(http.Header)
		if **in != nil {
			in, out := *in, *out
			*out = make(map[string][]string, len(*in))
			for key, val := range *in {
				var outVal []string
				if val == nil {
					(*out)[key] = nil
				} else {
					in, out := &val, &outVal
					*out = make([]string, len(*in))
					copy(*out, *in)
				}
				(*out)[key] = outVal
			}
		}
	}
	if in.Order != nil {
		in, out := &in.Order, &out.Order
		*out = new(int)
		**out = **in
	}
	if in.GetEndpointSuffix != nil {
		in, out := &in.GetEndpointSuffix, &out.GetEndpointSuffix
		*out = new(string)
		**out = **in
	}
	if in.PostEndpointSuffix != nil {
		in, out := &in.PostEndpointSuffix, &out.PostEndpointSuffix
		*out = new(string)
		**out = **in
	}
	if in.PatchEndpointSuffix != nil {
		in, out := &in.PatchEndpointSuffix, &out.PatchEndpointSuffix
		*out = new(string)
		**out = **in
	}
	if in.PutEndpointSuffix != nil {
		in, out := &in.PutEndpointSuffix, &out.PutEndpointSuffix
		*out = new(string)
		**out = **in
	}
	if in.DeleteEndpointSuffix != nil {
		in, out := &in.DeleteEndpointSuffix, &out.DeleteEndpointSuffix
		*out = new(string)
		**out = **in
	}
	if in.RuntimeVars != nil {
		in, out := &in.RuntimeVars, &out.RuntimeVars
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.HttpHeaders != nil {
		in, out := &in.HttpHeaders, &out.HttpHeaders
		*out = new(map[string]string)
		if **in != nil {
			in, out := *in, *out
			*out = make(map[string]string, len(*in))
			for key, val := range *in {
				(*out)[key] = val
			}
		}
	}
	in.Variables.DeepCopyInto(&out.Variables)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Action.
func (in *Action) DeepCopy() *Action {
	if in == nil {
		return nil
	}
	out := new(Action)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AuthConfig) DeepCopyInto(out *AuthConfig) {
	*out = *in
	if in.OAuth != nil {
		in, out := &in.OAuth, &out.OAuth
		*out = new(ConfigOAuth)
		(*in).DeepCopyInto(*out)
	}
	if in.CustomToken != nil {
		in, out := &in.CustomToken, &out.CustomToken
		*out = new(CustomToken)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AuthConfig.
func (in *AuthConfig) DeepCopy() *AuthConfig {
	if in == nil {
		return nil
	}
	out := new(AuthConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in AuthMap) DeepCopyInto(out *AuthMap) {
	{
		in := &in
		*out = make(AuthMap, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
		return
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AuthMap.
func (in AuthMap) DeepCopy() AuthMap {
	if in == nil {
		return nil
	}
	out := new(AuthMap)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConfigOAuth) DeepCopyInto(out *ConfigOAuth) {
	*out = *in
	if in.Scopes != nil {
		in, out := &in.Scopes, &out.Scopes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.EndpointParams != nil {
		in, out := &in.EndpointParams, &out.EndpointParams
		*out = make(map[string][]string, len(*in))
		for key, val := range *in {
			var outVal []string
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = make([]string, len(*in))
				copy(*out, *in)
			}
			(*out)[key] = outVal
		}
	}
	if in.ResourceOwnerUser != nil {
		in, out := &in.ResourceOwnerUser, &out.ResourceOwnerUser
		*out = new(string)
		**out = **in
	}
	if in.ResourceOwnerPassword != nil {
		in, out := &in.ResourceOwnerPassword, &out.ResourceOwnerPassword
		*out = new(string)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConfigOAuth.
func (in *ConfigOAuth) DeepCopy() *ConfigOAuth {
	if in == nil {
		return nil
	}
	out := new(ConfigOAuth)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CustomToken) DeepCopyInto(out *CustomToken) {
	*out = *in
	in.CustomAuthMap.DeepCopyInto(&out.CustomAuthMap)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CustomToken.
func (in *CustomToken) DeepCopy() *CustomToken {
	if in == nil {
		return nil
	}
	out := new(CustomToken)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in Seeders) DeepCopyInto(out *Seeders) {
	{
		in := &in
		*out = make(Seeders, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
		return
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Seeders.
func (in Seeders) DeepCopy() Seeders {
	if in == nil {
		return nil
	}
	out := new(Seeders)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StrategyConfig) DeepCopyInto(out *StrategyConfig) {
	*out = *in
	if in.AuthConfig != nil {
		in, out := &in.AuthConfig, &out.AuthConfig
		*out = make(AuthMap, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
	if in.Seeders != nil {
		in, out := &in.Seeders, &out.Seeders
		*out = make(Seeders, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StrategyConfig.
func (in *StrategyConfig) DeepCopy() *StrategyConfig {
	if in == nil {
		return nil
	}
	out := new(StrategyConfig)
	in.DeepCopyInto(out)
	return out
}
