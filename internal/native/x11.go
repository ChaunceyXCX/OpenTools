//go:build linux

package native

/*
#cgo linux pkg-config: gtk+-3.0
#cgo LDFLAGS: -lX11

#include <gtk/gtk.h>
#include <X11/Xlib.h>
#include <X11/Xatom.h>
#include <gdk/gdkx.h>

void set_popup_hint(GtkWindow *win) {
	gtk_window_set_type_hint(win, GDK_WINDOW_TYPE_HINT_POPUP_MENU);
}

void set_skip_taskbar(GtkWindow *win) {
	gtk_window_set_skip_taskbar_hint(win, TRUE);
	gtk_window_set_skip_pager_hint(win, TRUE);
}

void set_window_floating(GtkWindow *win) {
	// Get the X11 display and window
	GdkDisplay *display = gtk_widget_get_display(GTK_WIDGET(win));
	GdkWindow *gdk_window = gtk_widget_get_window(GTK_WIDGET(win));
	if (gdk_window == NULL) return;

	Display *xdisplay = gdk_x11_display_get_xdisplay(display);
	Window xwindow = gdk_x11_window_get_xid(gdk_window);

	// Set _NET_WM_WINDOW_TYPE to _NET_WM_WINDOW_TYPE_POPUP_MENU
	// This tells tiling WMs to float this window
	Atom net_wm_window_type = XInternAtom(xdisplay, "_NET_WM_WINDOW_TYPE", False);
	Atom net_wm_window_type_popup = XInternAtom(xdisplay, "_NET_WM_WINDOW_TYPE_POPUP_MENU", False);
	XChangeProperty(xdisplay, xwindow, net_wm_window_type, XA_ATOM, 32,
		PropModeReplace, (unsigned char *)&net_wm_window_type_popup, 1);

	// Set _NET_WM_STATE: SKIP_TASKBAR | SKIP_PAGER | ABOVE
	Atom net_wm_state = XInternAtom(xdisplay, "_NET_WM_STATE", False);
	Atom skip_taskbar = XInternAtom(xdisplay, "_NET_WM_STATE_SKIP_TASKBAR", False);
	Atom skip_pager = XInternAtom(xdisplay, "_NET_WM_STATE_SKIP_PAGER", False);
	Atom above = XInternAtom(xdisplay, "_NET_WM_STATE_ABOVE", False);
	Atom atoms[] = {skip_taskbar, skip_pager, above};
	XChangeProperty(xdisplay, xwindow, net_wm_state, XA_ATOM, 32,
		PropModeReplace, (unsigned char *)atoms, 3);

	XFlush(xdisplay);
}
*/
import "C"
import "unsafe"

func SetFloatingWindow(nativeWindow unsafe.Pointer) {
	if nativeWindow == nil {
		return
	}
	C.set_window_floating((*C.GtkWindow)(nativeWindow))
}
