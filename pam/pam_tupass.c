/*
 * TUPass PAM module
 */

/* 
 * Written by Cristian Gafton <gafton@redhat.com> 1996/09/10
 * Long password support by Philip W. Dalrymple <pwd@mdtsoft.com> 1997/07/18
 * Modified to use TUPass strength check by Simon Althaus and Alexander Peerdeman 2019/01/28
 * See the end of the file for original Copyright Information
 */

#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>
#include <string.h>
#include <syslog.h>
#include <stdarg.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <ctype.h>
#include <limits.h>
#include <pwd.h>
#include <security/pam_modutil.h>
#include "libtupass.h"

/*
 * here, we make a definition for the externally accessible function
 * in this file (this definition is required for static a module
 * but strongly encouraged generally) it is used to instruct the
 * modules include file to define the function prototypes.
 */

#define PAM_SM_PASSWORD

#include <security/pam_modules.h>
#include <security/_pam_macros.h>
#include <security/pam_ext.h>
#include <security/pam_appl.h>

/* argument parsing */
#define PAM_DEBUG_ARG 0x0001

struct tupass_options
{
    int threshold;
};

#define THRESHOLD 60

static int _pam_parse(pam_handle_t *pamh, struct tupass_options *opt,
           int argc, const char **argv)
{
    int ctrl = 0;

    /* step through arguments */
    for (ctrl = 0; argc-- > 0; ++argv)
    {
        char *ep = NULL;

        /* generic options */

        if (!strcmp(*argv, "debug"))
            ctrl |= PAM_DEBUG_ARG;
        else if (!strncmp(*argv, "type=", 5))
            pam_set_item(pamh, PAM_AUTHTOK_TYPE, *argv + 5);
        else if (!strncmp(*argv, "threshold=", 10))
        {
            opt->threshold = strtol(*argv + 10, &ep, 10);
            if (!ep || (opt->threshold < 0) || (opt->threshold > 100))
                opt->threshold = THRESHOLD;
        }
        else
        {
            pam_syslog(pamh, LOG_ERR, "pam_parse: unknown option; %s", *argv);
        }
    }

    return ctrl;
}

static const char *password_check(pam_handle_t *pamh, struct tupass_options *opt,
                                  const char *old, const char *new,
                                  const char *user)
{
    const char *msg = NULL;
    char *newmono = NULL;

    newmono = strdup(new);
    if (!newmono)
		msg = "memory allocation error";
    
    GoFloat64 result = CalculateStrength(newmono);

    if(result < opt->threshold) {
        msg = "TUPass fail";
        pam_error(pamh, "TUPass check failed (total strength: %f%%)", result);
    }

    if (newmono) {
		memset(newmono, 0, strlen(newmono));
		free(newmono);
	}

    return msg;
}

static int _pam_unix_approve_pass(pam_handle_t *pamh,
                                  unsigned int ctrl,
                                  struct tupass_options *opt,
                                  const char *pass_old,
                                  const char *pass_new)
{
    const char *msg = NULL;
    const char *user;
    int retval;

    if (pass_new == NULL)
    {
        if (ctrl & PAM_DEBUG_ARG)
            pam_syslog(pamh, LOG_DEBUG, "bad authentication token");
        pam_error(pamh, "%s", "No password has been supplied.");
        return PAM_AUTHTOK_ERR;
    }

    msg = password_check(pamh, opt, pass_old, pass_new, user);

    if (msg)
    {
        if (ctrl & PAM_DEBUG_ARG)
            pam_syslog(pamh, LOG_NOTICE,
                       "new passwd fails strength check: %s", msg);
        return PAM_AUTHTOK_ERR;
    };
    return PAM_SUCCESS;
}

/* The Main Thing (by Cristian Gafton, CEO at this module :-)
 * (stolen from http://home.netscape.com)
 */
int pam_sm_chauthtok(pam_handle_t *pamh, int flags, int argc, const char **argv)
{
    unsigned int ctrl;
    struct tupass_options options;

    memset(&options, 0, sizeof(options));
    options.threshold = THRESHOLD;

    ctrl = _pam_parse(pamh, &options, argc, argv);

    if (flags & PAM_PRELIM_CHECK)
    {
        return PAM_SUCCESS;
    }
    else if (flags & PAM_UPDATE_AUTHTOK)
    {
        int retval;
        const void *oldtoken;


        retval = pam_get_item(pamh, PAM_OLDAUTHTOK, &oldtoken);
        if (retval != PAM_SUCCESS)
        {
            if (ctrl & PAM_DEBUG_ARG)
                pam_syslog(pamh, LOG_ERR, "Can not get old passwd");
            oldtoken = NULL;
        }

        const char *crack_msg;
        const char *newtoken = NULL;

        retval = pam_get_authtok_noverify(pamh, &newtoken, NULL);
        if (retval != PAM_SUCCESS)
        {
            pam_syslog(pamh, LOG_ERR, "pam_get_authtok_noverify returned error: %s",
                        pam_strerror(pamh, retval));
            return retval;
        }
        else if (newtoken == NULL)
        { /* user aborted password change, quit */
            return PAM_AUTHTOK_ERR;
        }

        /* check it for strength... */
        retval = _pam_unix_approve_pass(pamh, ctrl, &options,
                                        oldtoken, newtoken);
        
        /* discard result as we only want to display warning */

        return PAM_SUCCESS;
    }
    else
    {
        if (ctrl & PAM_DEBUG_ARG)
            pam_syslog(pamh, LOG_NOTICE, "UNKNOWN flags setting %02X", flags);
        return PAM_SERVICE_ERR;
    }

    /* Not reached */
    return PAM_SERVICE_ERR;
}

/*
 * Copyright (c) Cristian Gafton <gafton@redhat.com>, 1996.
 *                                              All rights reserved
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 * 1. Redistributions of source code must retain the above copyright
 *    notice, and the entire permission notice in its entirety,
 *    including the disclaimer of warranties.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 * 3. The name of the author may not be used to endorse or promote
 *    products derived from this software without specific prior
 *    written permission.
 *
 * ALTERNATIVELY, this product may be distributed under the terms of
 * the GNU Public License, in which case the provisions of the GPL are
 * required INSTEAD OF the above restrictions.  (This clause is
 * necessary due to a potential bad interaction between the GPL and
 * the restrictions contained in a BSD-style copyright.)
 *
 * THIS SOFTWARE IS PROVIDED `AS IS'' AND ANY EXPRESS OR IMPLIED
 * WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES
 * OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED.  IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT,
 * INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT,
 * STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED
 * OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 * The following copyright was appended for the long password support
 * added with the libpam 0.58 release:
 *
 * Modificaton Copyright (c) Philip W. Dalrymple III <pwd@mdtsoft.com>
 *       1997. All rights reserved
 *
 * THE MODIFICATION THAT PROVIDES SUPPORT FOR LONG PASSWORD TYPE CHECKING TO
 * THIS SOFTWARE IS PROVIDED `AS IS'' AND ANY EXPRESS OR IMPLIED
 * WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES
 * OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED.  IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT,
 * INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT,
 * STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED
 * OF THE POSSIBILITY OF SUCH DAMAGE.
 * 
 * The following copyright was appended for inclusion of TUPass strength
 * check:
 * 
 * Modificaton Copyright (c) Simon Althaus <simon@alt.haus>
 *                           Alexander Peerdeman <a@peerdeman.info>
 *       2019. All rights reserved
 *
 * THE MODIFICATION THAT PROVIDES SUPPORT FOR LONG PASSWORD TYPE CHECKING TO
 * THIS SOFTWARE IS PROVIDED `AS IS'' AND ANY EXPRESS OR IMPLIED
 * WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES
 * OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED.  IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT,
 * INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT,
 * STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED
 * OF THE POSSIBILITY OF SUCH DAMAGE.
 * 
 */
